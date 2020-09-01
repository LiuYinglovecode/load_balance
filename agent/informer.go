package agent

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"time"

	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	cron "github.com/robfig/cron"
)

// InfomerKeyPrefx agent 定时汇报到etcd中的key前缀
// 存在etcd中定时汇报的路径是/alb/agent/{agentid}
const InfomerKeyPrefx = "/alb/agent/"

// Informer 负责向etcd汇报自己的状态
type Informer interface {
	Start() error
	Stop() error
	HeartBeat() error
}

type etcdInformer struct {
	AgentID     int64
	Role        string
	RPC         string
	reportTimer *cron.Cron
	cli         *clientv3.Client
	store       LBPolicyStore
}

// NewInformer create informer
func NewInformer(config *Config, store LBPolicyStore) (Informer, error) {
	timer := cron.New()
	endpoints := config.Endpoints

	var tlsconf *tls.Config
	tlsconf, err := tlsConfig(config.EtcdCAPath, config.EtcdCertPath, config.EtcdKeyPath)
	if err != nil {
		sysLogger.Errorf("create tls config failed, reason: %v", err)
	}

	etcdcli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		TLS:         tlsconf,
	})
	if err != nil {
		sysLogger.Errorf("Fail to connect etcd endpoints: %+v, reason: %v", endpoints, err)
		return nil, err
	}

	// 如果非高可用部署，则ROLE定义为MASTER
	// 否则读入配置文件中定义的STATE
	var role string
	if config.KeepalivedCfg == nil {
		role = "MASTER"
	} else {
		role = config.KeepalivedCfg.State
	}

	informer := &etcdInformer{
		config.AgentID,
		role,
		config.RPC,
		timer,
		etcdcli,
		store,
	}

	err = informer.reportTimer.AddFunc("@every 5m", func() {
		sysLogger.Debug("Report agent status to etcd every 5 minitue ...")
		err := informer.HeartBeat()
		if err != nil {
			sysLogger.Errorf("Report agent status to etcd failed, reason: %v", err)
		}
	})
	if err != nil {
		sysLogger.Errorf("Add timer function into informer failed, reason: %v", err)
		return nil, err
	}
	return informer, nil
}

// Start 启动定时汇报
func (e *etcdInformer) Start() error {
	e.HeartBeat()
	e.reportTimer.Start()
	return nil
}

// Stop 停止定时汇报
func (e *etcdInformer) Stop() error {
	e.reportTimer.Stop()
	return nil
}

// HeartBeat 获取agent状态数据， 上传到etcd中
func (e *etcdInformer) HeartBeat() error {
	hostname, err := os.Hostname()
	if err != nil {
		sysLogger.Errorf("Get Hostname failed, reason: %v", err)
	}
	hostip, err := getHostIP()
	if err != nil {
		sysLogger.Errorf("Get Hostip failed, reason: %v", err)
	}

	policies := e.store.Get()

	informData := &model.LBAgentStauts{
		ID:            e.AgentID,
		Role:          e.Role,
		ControllerRPC: e.RPC,
		API:           "",

		HostIP:   hostip,
		HostName: hostname,

		Policies: policies,
		TimeAt:   model.NewADCTime(time.Now()),
	}

	strID := strconv.FormatInt(e.AgentID, 10)

	value, err := json.Marshal(informData)
	if err != nil {
		sysLogger.Errorf("Informer data json marshall failed, reason: %v", err)
		return err
	}

	_, err = e.cli.Put(context.TODO(), InfomerKeyPrefx+strID+"/"+e.Role, string(value))
	if err != nil {
		sysLogger.Errorf("Put informer data into etcd failed, reason: %v", err)
		return err
	}
	return nil
}

// EtcdClient 返回etcd 客户端
func (e *etcdInformer) EtcdClient() *clientv3.Client {
	return e.cli
}

// getHostIP 遍历本机上的网卡，返回ip地址的数组的json字符串,
func getHostIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ips []string

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip.String())
		}
	}

	ret, err := json.Marshal(ips)
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

// 配置TLS
func tlsConfig(ca, cert, key string) (*tls.Config, error) {
	var cfgtls *transport.TLSInfo
	tlsinfo := transport.TLSInfo{}

	tlsinfo.CertFile = cert
	tlsinfo.KeyFile = key
	tlsinfo.TrustedCAFile = ca
	cfgtls = &tlsinfo

	return cfgtls.ClientConfig()
}
