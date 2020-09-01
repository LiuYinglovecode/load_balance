package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.htres.cn/casicloud/alb/apis"
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/center/queue"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/pkg/errors"
)

// LBRequestHandler 处理请求的业务逻辑接口
type LBRequestHandler interface {
	LBRequestAcceptor
	LBRequestProcessor
}

// LBRequestAcceptor 记录请求接口
type LBRequestAcceptor interface {
	// 接受请求消息, 记录请求,
	// 返回request id
	AcceptRequest(request *model.LBRequest) (int64, error)
	// 查询请求处理状态, 返回状态码
	QueryRequest(id int64) (int32, error)
}

// LBRequestProcessor 处理请求接口
type LBRequestProcessor interface {
	// 按顺序处理已经受理的消息, 建立代理通道
	// 将处理后消息记录到持久化存储中
	ProcessRequest()
}

// RequestHandlerImpl 实现客户端请求消息的处理
type RequestHandlerImpl struct {
	lbrRepo         dao.LBRecordDao
	lbpRepo         dao.LBPoolDao
	lbreqRepo       dao.LBRequestRepository
	queueController queue.MessageQueueHandler
	agentClient     Communicator
}

// NewRequestHandlerImpl 构造函数
func NewRequestHandlerImpl(
	lbrRepo dao.LBRecordDao,
	lbpRepo dao.LBPoolDao,
	lbreqRepo dao.LBRequestRepository,
	que queue.MessageQueueHandler,
	agentCli Communicator) LBRequestHandler {
	return &RequestHandlerImpl{lbrRepo, lbpRepo, lbreqRepo, que, agentCli}
}

// AcceptRequest 接收客户端请求，将请求加入队列，并写入数据库
func (r *RequestHandlerImpl) AcceptRequest(request *model.LBRequest) (int64, error) {
	// 从请求中获取用户信息
	err := r.extractUserInfo(request)
	if err != nil {
		common.SysLogger.Errorf("extract user info from lbrequest failedd, request: %+v, reason: %v", *request, err)
		return 0, err
	}

	// 更具请求获取处理该请求的agentid
	agent, err := r.getLbAgentID(request)
	if err != nil {
		common.SysLogger.Errorf("Get agent id form lbrequest failed, request: %+v, reason: %v", *request, err)
		return 0, err
	}

	// 先写数据库
	// todo: 增加事务处理
	err = r.lbreqRepo.Create(request)
	if err != nil {
		common.SysLogger.Errorf("LBRequest persist into db failed, request: %+v, reason: %v", *request, err)
		return 0, err
	}

	// 再写消息队列
	if err := r.queueController.Enqueue(strconv.FormatInt(agent, 10), request); err != nil {
		common.SysLogger.Errorf("LBRequest enqueue failed, request: %+v, reason: %v", *request, err)
		return 0, err
	}
	common.SysLogger.Infof("LBRequest accepted, request: %+v", *request)
	return request.RequestID, nil
}

// QueryRequest Client根据request id查询请求状态
func (r *RequestHandlerImpl) QueryRequest(id int64) (int32, error) {
	req, err := r.lbreqRepo.GetByID(id)
	if err != nil {
		common.SysLogger.Errorf("LBRequest query failed, request id: %d, reason: %v", id, err)
		return model.StatusUnknown, err
	}
	return req.Status, nil
}

// ProcessRequest 处理请求
// 该方法无限循环，应该使用单独的线程启动
func (r *RequestHandlerImpl) ProcessRequest() {
	for {
		req, agentID := r.queueController.WatchAndDequeue()
		common.SysLogger.Infof("Begin process request: %+v", *req)
		if err := r.doProcessRequest(req, agentID); err != nil {
			common.SysLogger.Errorf("Process request faild, request %+v, reason %v", *req, err)
		}
	}
}

// doProcessRequest 进行如下处理
// 写数据库
// 通过grpc与agent通信
// 如果有处理失败，则将lbrequest在数据库中标记为处理失败
func (r *RequestHandlerImpl) doProcessRequest(request *model.LBRequest, agentID string) error {
	request.Status = model.StatusInProcessing
	// 将request标记为正在处理
	if err := r.lbreqRepo.Update(request); err != nil {
		if err1 := r.handleAgentCommunicateFail(request); err1 != nil {
			return errors.Wrap(err, err1.Error())
		}
		return err
	}

	// 从agent的状态汇报中获取agent地址
	rpcAddrs, err := r.getLbAgentAddrFromEtcd(agentID)
	if err != nil {
		if err1 := r.handleAgentCommunicateFail(request); err1 != nil {
			return errors.Wrap(err, err1.Error())
		}
		return err
	}

	// 与主备 agent 分别通信
	errs := make([]error, len(rpcAddrs))
	for i, rpcAddr := range rpcAddrs {
		errs[i] = r.conmunicateWithAgent(rpcAddr, request)
	}

	if gotNoSuccess(errs) {
		err := errors.New("conmmunicate with all the agents failed" + errs[0].Error())
		if err1 := r.handleAgentCommunicateFail(request); err1 != nil {
			return errors.Wrap(err, err1.Error())
		}
		return err
	}

	return nil
}

// gotNoSuccess 检查通信成功的情况, 如果全部失败，返回true
func gotNoSuccess(errs []error) bool {
	for _, e := range errs {
		if e == nil {
			return false
		}
		common.SysLogger.Errorf("communicate with agent failed, reason: %s", e)
	}
	return true
}

// conmunicateWithAgent 与agent通信，并将结果写入数据库
func (r *RequestHandlerImpl) conmunicateWithAgent(agentAddr string, request *model.LBRequest) error {
	result, err := r.agentClient.SendLBRequest(agentAddr, request)
	if err != nil {
		return err
	}

	if result.Code == apis.RetOk {
		common.SysLogger.Infof("Agent rpc call success %+v", *result)
		request.Status = model.StatusHandleSuccess

		// 更新request
		if err := r.lbreqRepo.Update(request); err != nil {
			return err
		}

		// 更新lbr状态
		lbr := request.Policy.Record

		attrs := make(map[string]interface{})

		switch request.Action {
		case model.ActionAdd:
			attrs["status"] = model.InUse
		case model.ActionStop:
			attrs["status"] = model.Applied
		}

		if _, err := r.lbrRepo.UpdateAttribute(lbr.ID, attrs); err != nil {
			return err
		}
	} else {
		common.SysLogger.Warnf("Agent rpc call return %+v", *result)
		return fmt.Errorf("Agent return error: %v", result.Msg)
	}
	return nil
}

// 如果agent通信失败，进行错误处理
func (r *RequestHandlerImpl) handleAgentCommunicateFail(request *model.LBRequest) error {
	request.Status = model.StatusHandleFailed
	if err := r.lbreqRepo.Update(request); err != nil {
		return err
	}
	return nil
}

// reEnqueue 将处理失败的请求重新加入队列 目前未使用
func (r *RequestHandlerImpl) reEnqueue(request *model.LBRequest) error {
	agentID, err := r.getLbAgentID(request)
	if err != nil {
		common.SysLogger.Errorf("Get agent id form lbrequest failed, request: %+v, reason: %v", *request, err)
		return err
	}

	if err := r.queueController.Enqueue(strconv.FormatInt(agentID, 10), request); err != nil {
		common.SysLogger.Errorf("LBRequest enqueue failed after process failure, request: %+v, reason: %v", *request, err)
		return err
	}
	return nil
}

// getLbAgentId 返回处理该请求的agent的id值
// 处理流程如下:
// 1. 根据ip或者域名查lbr是否存在, 根据ip确定由哪个agent处理, 得到agent id
// 2. todo: 如果ip为空那么看该用户是否有权限从ip池分配ip,
//       如果有，分配的ip加入lbr, ip值加入到lbrequest
func (r *RequestHandlerImpl) getLbAgentID(request *model.LBRequest) (int64, error) {
	conditions := make(map[string]interface{})

	user := request.Policy.Record.Owner.String()
	conditions["owner"] = user

	switch request.Policy.Record.Type {
	case model.TypeIP:
		conditions["ip"] = request.Policy.Record.IP.String()
		conditions["port"] = request.Policy.Record.Port
	case model.TypeDomain:
		conditions["domain"] = request.Policy.Record.Domain.String()
	}

	return r.doGetLBAgentID(conditions, request.Policy.Record.Type, request)
}

// 用于验证k8s namespace名称中是否包含用户信息的正则
const namespaceRegex = "^[a-z]+-[0-9]{10,}$"

// extractUserInfo 如果LBPolicy中不存在User信息，那么从LBRequest信息中提取UserID
func (r *RequestHandlerImpl) extractUserInfo(request *model.LBRequest) error {
	if len(request.Policy.Record.Owner.String()) != 0 {
		return nil
	}

	reg, err := regexp.Compile(namespaceRegex)
	if err != nil {
		common.SysLogger.Errorf("creat namespace regex faild, reason: %s", err)
		return err
	}

	// 如果namespace名称包含userid, 使用此id
	if reg.MatchString(request.User.String()) {
		nameAndUser := strings.Split(request.User.String(), "-")
		request.Policy.Record.Owner = model.NewADCString(nameAndUser[1])
	} else {
		// 如果namespace名称不包含userid, 使用此namespace
		// 这种情况主要为了兼容非adc平台创建的namespace
		request.Policy.Record.Owner = request.User
	}
	return nil
}

// todo: 逻辑过于复杂需要重构
func (r *RequestHandlerImpl) doGetLBAgentID(conditions map[string]interface{}, t model.LBType, request *model.LBRequest) (int64, error) {
	record := fmt.Sprintf("lbr: %+v", conditions)

	lbrs := []model.LBRecord{}
	err := r.lbrRepo.ListLBRecord(conditions, &lbrs)
	if err != nil {
		common.SysLogger.Errorf("Check lbr failed, %s, reason: %s", record, err)
		return 0, err
	}

	if len(lbrs) == 0 {
		common.SysLogger.Infof("Lbr not existed, %s", record)
		return 0, fmt.Errorf("Lbr not existed, %s", record)
	}

	// 查询出多个lbr的情况
	// 这种情况应该在添加lbr时候进行处理
	// 如果出现了，不处理，返回错误信息
	if len(lbrs) > 1 {
		msg := "get multiple lbr from database, Can not go on, %+v"
		common.SysLogger.Errorf(msg, lbrs)
		return 0, fmt.Errorf(msg, lbrs)
	}

	// 该lbr是否已经在使用
	// todo: 建立代理后更新status
	lbr := lbrs[0]

	if lbr.Status == model.InUse && request.Action == model.ActionAdd {
		msg := "lbr is already in use, lbr: %s"
		common.SysLogger.Errorf(msg, record)
		return 0, fmt.Errorf(msg, record)
	}

	// 填充lbr id到 lbrequest中, lbreqesut处理完成后, 更新lbr状态需要使用
	request.Policy.Record.ID = lbr.ID

	// 获取lbpool
	ll := []model.LBPool{}
	err = r.lbpRepo.ListLBPool(make(map[string]interface{}), &ll)
	if err != nil {
		common.SysLogger.Errorf("fail to list lbpool, reason: err")
		return 0, nil
	}

	// 选择出处理该lbr的lbpool
	lbpools := make([]model.LBPool, 0)
	if t == model.TypeIP {
		lbpools = filterLBPool(ll, func(pool model.LBPool) bool {
			// todo: 将ip的字符串类型统一
			if lbr.IP.String() == pool.IP.String() && lbr.Port >= pool.StartPort && lbr.Port <= pool.EndPort {
				return true
			}
			return false
		})
	} else {
		lbpools = filterLBPool(ll, func(pool model.LBPool) bool {
			r, err := regexp.Compile(pool.DomainRegex.String())
			common.SysLogger.Errorf("lbpool domain regrex complie failed, reason: %v", err)
			if r.MatchString(lbr.Domain.String()) {
				return true
			}
			return false
		})
	}

	if len(lbpools) == 0 {
		common.SysLogger.Errorf("no lbpool to handle lbr: %s", record)
		return 0, fmt.Errorf("no lbpool to handle lbr: %s", record)
	}

	// 正常情况一个lbr只属于一个资源池
	if len(lbpools) > 1 {
		common.SysLogger.Warnf("found multiple lbpools to handle lbr: %s", record)
	}

	// 随机分配一个agent处理请求
	rand.Seed(time.Now().UnixNano())

	numOfAgents := len(lbpools[0].Agents)
	if numOfAgents == 0 {
		msg := "found no agent to handle lbpool, lbp: %+v"
		common.SysLogger.Warnf(msg, lbpools[0])
		return 0, fmt.Errorf(msg, lbpools[0])
	}
	randoAgentID := rand.Intn(numOfAgents)

	return lbpools[0].Agents[randoAgentID].ID, nil
}

// filterLBPool 根据条件过滤lbpool数组
func filterLBPool(pools []model.LBPool, f func(pool model.LBPool) bool) []model.LBPool {
	vsf := make([]model.LBPool, 0)
	for _, v := range pools {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// getLbAgentAddr 根据AgentID读取etcd, 返回agent rpc地址, 包括master和backup
func (r *RequestHandlerImpl) getLbAgentAddrFromEtcd(agentID string) ([]string, error) {
	etcdcli := r.queueController.(*queue.EtcdReadWriter)
	values, err := etcdcli.EtcdGetWithPrefix(common.InfomerKeyPrefx + agentID)
	if err != nil {
		return nil, err
	}

	var status []string
	for _, v := range values {
		s := model.LBAgentStauts{}

		err = json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.SysLogger.Errorf("lbagent status json unmarshll failed, reason: %v", err)
			return nil, err
		}

		status = append(status, s.ControllerRPC)
	}

	return status, nil
}

// getIPFromPool 从IP池中分配IP
// todo: 暂时未实现
func (r *RequestHandlerImpl) getIPFromPool(request *model.LBRequest) (string, error) {
	panic("implement me")
}
