package cloudprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"code.htres.cn/casicloud/alb/pkg/client"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	log "github.com/sirupsen/logrus"
)

const (
	dialTimeout     = 5 * time.Second
	lbCreateTimeout = 100 * time.Second
)

// Config 配置文件
type Config struct {
	Global struct {
		ApiCAFile  string `gcfg:"apica"`
		EtcdCAFile string `gcfg:"etcdca"`
		Keyfile    string `gcfg:"keyfile"`
		Certfile   string `gcfg:"certfile"`
		//LBKeyPrefix   etcd 中存储的负载均衡的key
		LBKeyPrefix string `gcfg:"lbkeyprefix"`
		// Endpoints etcd 地址
		Endpoints []string `gcfg:"endpoints"`
		// LbmcUrl center 地址
		LbmcUrl   string `gcfg:"lbmcurl"`
		APIServer string `gcfg:"apiserver"`
	}
}

// LBController alb控制器
type LBController struct {
	Conf   *Config
	Lbmcli *client.AlbClient
}

// NewAlbCli 构造函数
func NewAlbCli(conf *Config, cli *client.AlbClient) *LBController {
	if cli == nil {
		cli = client.NewAlbClient(conf.Global.LbmcUrl, nil)
	}

	return &LBController{
		conf,
		cli,
	}
}

func (a *LBController) newEtcdClientCfg() (*clientv3.Config, error) {
	var cfgtls *transport.TLSInfo
	tlsinfo := transport.TLSInfo{}

	cfg := &clientv3.Config{
		Endpoints:   a.Conf.Global.Endpoints,
		DialTimeout: dialTimeout,
	}

	tlsinfo.CertFile = a.Conf.Global.Certfile
	tlsinfo.KeyFile = a.Conf.Global.Keyfile
	tlsinfo.TrustedCAFile = a.Conf.Global.EtcdCAFile
	cfgtls = &tlsinfo

	clientTLS, err := cfgtls.ClientConfig()
	if err != nil {
		return nil, err
	}
	cfg.TLS = clientTLS

	return cfg, nil
}

func (a *LBController) newEtcdClient() (*clientv3.Client, error) {
	cfg, err := a.newEtcdClientCfg()
	if err != nil {
		log.WithField("reason", err).Error("Get Etcd Client config failed")
		return nil, err
	}

	etcdcli, err := clientv3.New(*cfg)
	if err != nil {
		log.WithField("reason", err).Error("Get Etcd Client failed")
		return nil, err
	}

	return etcdcli, nil
}

//GetLoadBalancerFromEtcd retrieve result from etcd
func (a *LBController) GetLoadBalancerFromEtcd(namespace string, serviceName string) ([]model.LBPolicy, error) {
	lbkey := fmt.Sprintf("%s/%s/%s", a.Conf.Global.LBKeyPrefix, namespace, serviceName)
	text, err := a.etcdGet(lbkey)
	if err != nil {
		return nil, err
	}

	st := []model.LBPolicy{}
	err = json.Unmarshal([]byte(text), &st)
	if err != nil {
		log.WithField("reason", err).Error("Lbpolicy[] json unmarshall failed")
		return nil, err
	}
	return st, nil
}

// CreateLoadBalancer 申请创建负载均衡
// 返回LBRequestID
func (a *LBController) CreateLoadBalancer(request *model.LBRequest) (string, error) {
	apiResult, err := a.Lbmcli.CreateLB(request)
	if err != nil {
		log.WithField("reason", err).Error("Send create lb request, but got err from lbmc")
		return "", err
	}

	if apiResult.Code != model.Ok {
		return "", fmt.Errorf(apiResult.Message)
	}

	if data, ok := apiResult.Data.(map[string]interface{}); ok {
		if id, okk := data["request_id"]; okk {
			return id.(string), nil
		}
	}
	return "", fmt.Errorf("CreateLB parse reqeust_id from lbmc response failed, raw respponse is %+v", *apiResult)
}

// DeleteLoadBalancer 发送删除负载均衡请求，返回request_id
func (a *LBController) DeleteLoadBalancer(request *model.LBRequest) (string, error) {
	apiResult, err := a.Lbmcli.StopLB(request)
	if err != nil {
		return "", err
	}

	if apiResult.Code != model.Ok {
		return "", fmt.Errorf(apiResult.Message)
	}

	if data, ok := apiResult.Data.(map[string]interface{}); ok {
		if id, okk := data["request_id"]; okk {
			return id.(string), nil
		}
	}
	return "", fmt.Errorf("StopLB parse reqeust_id from lbmc response failed, raw respponse is %+v", *apiResult)
}

//EtcdGet get
func (a *LBController) etcdGet(key string) (string, error) {
	cli, err := a.newEtcdClient()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	var resp *clientv3.GetResponse
	if resp, err = cli.Get(context.Background(), key); err != nil {
		log.WithField("reason", err).WithField("key", key).Error("Get etcd key failed")
		return "", err
	}
	log.WithField("resp", resp.Kvs).Info("Get key from etcd success")
	var ret string
	for _, kv := range resp.Kvs {
		ret += string(kv.Value)
	}
	return ret, nil
}

//EtcdPut put
func (a *LBController) etcdPut(key string, val string) error {
	cli, err := a.newEtcdClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	var resp *clientv3.PutResponse
	if resp, err = cli.Put(context.Background(), key, val); err != nil {
		log.WithField("reason", err).WithField("key", key).Error("Put etcd key field")
		return err
	}

	log.WithField("resp", resp).Info("Put key to etcd success")
	return nil
}

// WaitforLbReady wait for req finish
// 轮询lbmc
func (a *LBController) waitforLbReady(reqID string) error {
	for {
		id, err := strconv.ParseInt(reqID, 10, 64)
		if err != nil {
			log.WithField("reason", err).Errorf("Parse requestId failed")
			return err
		}
		log.WithField("reqID", reqID).Info("Waiting for Lb ready... ")

		apiResult, err := a.Lbmcli.QueryLB(id)

		if err != nil || apiResult.Code != model.Ok {
			// 网络请求失败，不返回err, 2秒后2秒后重试
			time.Sleep(2 * time.Second)
			log.WithField("reason", err).Error("Query LB state from lbmc failed, retry...")
			continue
		}

		// 返回数据解析失败 重试
		var status int32 = -1
		if data, ok := apiResult.Data.(map[string]interface{}); ok {
			if s, okk := data["status"]; okk {
				status = int32(s.(float64))
			}
		}

		if status == -1 {
			time.Sleep(2 * time.Second)
			log.WithField("resp", apiResult).Error("Query LB state from lbmc, parse resp failed, retry...")
			continue
		}

		// 如果状态是未处理或者正在处理，每2秒尝试获取一次状态
		// 处理成功返回nil
		// 处理失败返回error
		// todo: 不应该通过是否有error来确认是否返回成功
		switch int32(status) {
		case model.StatusUnHandle:
			time.Sleep(2 * time.Second)
			continue
		case model.StatusInProcessing:
			time.Sleep(2 * time.Second)
			continue
		case model.StatusHandleSuccess:
			log.WithField("requestID", reqID).Info("Success from lbmc")
			return nil
		default:
			return fmt.Errorf("LBMC return msg : create LB failed, reason: %+v", apiResult)
		}
	}
}

//CheckLbReady check lb existence
func (a *LBController) CheckLbReady(namespace string, service string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", a.Conf.Global.LBKeyPrefix, namespace, service)
	text, err := a.etcdGet(key)
	if err != nil {
		return "", err
	}

	st := []model.LBPolicy{}
	err = json.Unmarshal([]byte(text), &st)
	if err != nil || len(st) == 0 {
		log.WithField("reason", err).WithField("policies", st).Error("Lbpolicy[] json unmarshall failed")
		return "", errors.Wrap(err,"LB not ready")
	}
	return st[0].Record.IP.String(), nil
}

// PersistPolicies 申请IP成功后将policy 存储到etcd中
func (a *LBController) PersistPolicies(namespace string, service string, policies *[]model.LBPolicy) error {
	key := fmt.Sprintf("%s/%s/%s", a.Conf.Global.LBKeyPrefix, namespace, service)
	by, err := json.Marshal(policies)
	if err != nil {
		log.WithField("reason", err).Error("Marshall policies to json failed")
		return err
	}

	err = a.etcdPut(key, string(by))
	if err != nil {
		log.WithField("reason", err).Error("Put policies into etcd failed")
		return err
	}
	return nil
}

// DeletePolicies 删除etcd中存储的 policies
func (a *LBController) DeletePolicies(namespace string, service string) error {
	key := fmt.Sprintf("%s/%s/%s", a.Conf.Global.LBKeyPrefix, namespace, service)
	err := a.etcdDel(key)
	if err != nil {
		log.WithField("reason", err).Error("Delete policies from etcd failed")
		return err
	}
	return nil
}

func (a *LBController) etcdDel(key string) error {
	cli, err := a.newEtcdClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	var resp *clientv3.DeleteResponse
	if resp, err = cli.Delete(context.Background(), key); err != nil {
		log.WithField("reason", err).WithField("key", key).Error("del etcd key faild")
		return err
	}
	log.WithField("key", key).WithField("resp", resp).Info("del etcd key success")
	return nil
}
