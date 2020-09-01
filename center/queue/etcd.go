package queue

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/center/dao"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/coreos/etcd/clientv3"
	recipe "github.com/coreos/etcd/contrib/recipes"
	"github.com/coreos/etcd/pkg/transport"
)

// LBAgentEventType 表示client发送的请求事件类型
type LBAgentEventType int8

const (
	// LBAgentAdd 表示添加负载均衡事件类型
	LBAgentAdd LBAgentEventType = 0
	// LBAgentDelete 表示删除负载均衡事件类型
	LBAgentDelete LBAgentEventType = 1
)

// LBAgentEvent 表示Agent相关的事件
type LBAgentEvent struct {
	EventType LBAgentEventType
	AgentName string
	// 如果是agent删除事件，true表示同时清除消息队列中的LBRequest
	Purge bool
}

var lBAgentEventChannel = make(chan *LBAgentEvent)

// GetLBAgentEventChannel 用于监视LBAgent的改变，调整相应的队列
// LBAgent可能由运维人员手动添加
// 也可能由于节点异常下线
// 使用全局变量造成的问题是：状态存在单独channel中,无法多线程使用EtcdReadWriter
// 好处是：代码简单，如果由多线程需求后需修改
func GetLBAgentEventChannel() chan *LBAgentEvent {
	return lBAgentEventChannel
}

// EtcdReadWriter 读写etcd, 实现消息队列
type EtcdReadWriter struct {
	client *clientv3.Client
	queues map[string]*recipe.Queue

	lbaRepo dao.LBAgentRepository

	startedWatch   bool
	requestChannel chan queueWatchEvent
}

// NewEtcdMessageQueueHandler 构造函数
func NewEtcdMessageQueueHandler(endpoints []string, ca, cert, key string, repo dao.LBAgentRepository) (MessageQueueHandler, error) {
	if len(endpoints) == 0 {
		endpoints = append(endpoints, "localhost:2379")
	}

	var config *tls.Config
	config, err := tlsConfig(ca, cert, key)
	if err != nil {
		common.SysLogger.Errorf("create tls config failed, reason: %v", err)
	}
	// 创建etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		TLS:         config,
	})
	if err != nil {
		common.SysLogger.Errorf("Fail to connect etcd endpoints: %+v, reason: %v", endpoints, err)
		return nil, err
	}

	queues := make(map[string]*recipe.Queue)

	names, err := getQueueNames(repo)
	if err != nil {
		common.SysLogger.Errorf("Can not get agent name, reason: %v", err)
		return nil, err
	}

	for _, v := range names {
		k := common.LBReqKeyPrefix + v
		queues[k] = recipe.NewQueue(client, k)
	}

	reqChannel := make(chan queueWatchEvent)
	return &EtcdReadWriter{client, queues, repo, false, reqChannel}, nil
}

func tlsConfig(ca, cert, key string) (*tls.Config, error) {
	// 配置TLS
	var cfgtls *transport.TLSInfo
	tlsinfo := transport.TLSInfo{}

	tlsinfo.CertFile = cert
	tlsinfo.KeyFile = key
	tlsinfo.TrustedCAFile = ca
	cfgtls = &tlsinfo

	return cfgtls.ClientConfig()
}

// OnAgentChange 监听agent的变化, 修改对应的消息队列
// todo: 需要增加mux
func (e *EtcdReadWriter) OnAgentChange() {
	for event := range lBAgentEventChannel {
		keyName := common.LBReqKeyPrefix + event.AgentName
		switch event.EventType {
		case LBAgentAdd:
			// queueName已经存在
			if _, ok := e.queues[keyName]; ok {
				common.SysLogger.Infof("Received agent add event,but coresponding queue %s already exists", keyName)
				return
			}
			e.queues[keyName] = recipe.NewQueue(e.client, keyName)
			common.SysLogger.Infof("Handled LBAgentEvent %+v", event)

		case LBAgentDelete:
			if _, ok := e.queues[keyName]; ok {
				delete(e.queues, keyName)
				common.SysLogger.Infof("Deleted queue %s", keyName)
			}
			// 清除etcd中存储的LBRequest
			if event.Purge {
				_, err := e.client.Delete(context.TODO(), keyName, clientv3.WithPrefix())
				if err != nil {
					common.SysLogger.Errorf("Error purge queue %s, reason: %v", keyName, err)
				}
			}
			common.SysLogger.Infof("Handled LBAgentEvent %+v", event)
		default:
			common.SysLogger.Warnf("Received unknown LbAgentEvent %+v", event)
		}
	}
}

// WatchAndDequeue 为每一个队列启动一个线程, 线程间通过channel通信，
// 如果etcd中存在未处理的request, 那么返回该request和队列的名字(agent id）
func (e *EtcdReadWriter) WatchAndDequeue() (*model.LBRequest, string) {
	// 注意: 目前单线程调用WatchAndDequeue, 如果多线程调用需要给false加锁
	if e.startedWatch == false {
		for k := range e.queues {
			go e.doWatchAndDequeue(k, e.requestChannel)
		}
		e.startedWatch = true
	}
	common.SysLogger.Infof("EtcdReaderWriter is Waiting for LBRequest...")
	event := <-e.requestChannel

	p := strings.Split(event.queueName, "/")

	return event.request, p[len(p)-1]
}

type queueWatchEvent struct {
	request   *model.LBRequest
	queueName string
}

func (e *EtcdReadWriter) doWatchAndDequeue(queueName string, out chan<- queueWatchEvent) {
	for {
		queue, ok := e.queues[queueName]
		if !ok {
			common.SysLogger.Warnf("Queue %s id remove, watch process stop and exit", queueName)
			return
		}

		s, err := queue.Dequeue()
		if err != nil {
			common.SysLogger.Errorf(" Execute dequeue failed, queue: %+v, reason: %v", *queue, err)
		}
		common.SysLogger.Infof("dequeued LBRequest %v", s)

		lbrequest := &model.LBRequest{}

		err = json.Unmarshal([]byte(s), lbrequest)
		if err != nil {
			common.SysLogger.Errorf("LBRequest %s json unmarshal failed, reason %v", string(s), err)
		}
		event := queueWatchEvent{
			lbrequest,
			queueName,
		}

		out <- event
	}
}

// Enqueue 事件入列操作
func (e *EtcdReadWriter) Enqueue(queueName string, request *model.LBRequest) error {
	fullQueueName := common.LBReqKeyPrefix + queueName
	value, err := json.Marshal(request)
	if err != nil {
		common.SysLogger.Errorf("LBRequest %+v json unmarshal failed, reason %v", *request, err)
		return err
	}
	if queue, ok := e.queues[fullQueueName]; ok {
		err := queue.Enqueue(string(value))
		if err != nil {
			common.SysLogger.Errorf(" LBRequest enqueue failed, request: %v, reason: %v", string(value), err)
			return err
		}
	} else {
		common.SysLogger.Warnf("Received lbrequest %+v, but queue not existed", *request)
		return fmt.Errorf("queue not existed, quqeue_name(agent_id): %s", queueName)
	}
	common.SysLogger.Infof("Enqueued request %+v", *request)
	return nil
}

// GetCurrentQueueNames 返回当前全部队列名称
func (e *EtcdReadWriter) GetCurrentQueueNames() []string {
	names := make([]string, 0, len(e.queues)+1)
	for k := range e.queues {
		names = append(names, k)
	}
	common.SysLogger.Infof("Current queue names %v", names)
	return names
}

// getQueueNames 返回当前所有agent名称
// 读取数据库，获取所有的LBAgentId
func getQueueNames(repo dao.LBAgentRepository) ([]string, error) {
	var ids []string

	agents, err := repo.List(make(map[string]interface{}))
	if err != nil {
		common.SysLogger.Errorf("Get agentId from database failed, reason: %v", err)
		return nil, err
	}
	for _, v := range agents {
		ids = append(ids, strconv.FormatInt(v.ID, 10))
	}
	return ids, nil
}

// EtcdGetWithPrefix 从etcd中根据key前缀读取value，返回value数组
func (e *EtcdReadWriter) EtcdGetWithPrefix(key string) ([]string, error) {
	resp, err := e.client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		common.SysLogger.Errorf("get etcd key prefix: %s faild, reason: %s", key, err)
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		msg := "found no agent status in etcd,  key prefix: %s "
		common.SysLogger.Errorf(msg, key)
		return nil, fmt.Errorf(msg, key)
	}

	var values []string
	for _, ev := range resp.Kvs {
		values = append(values, string(ev.Value))
	}
	common.SysLogger.Infof("success get from etcd, key prefix: %v, values %v", key, values)

	return values, nil
}
