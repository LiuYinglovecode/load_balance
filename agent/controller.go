package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"code.htres.cn/casicloud/alb/pkg/model"
	dockercli "github.com/docker/docker/client"
)

// LBController loadbalance controller
// 其中一个LBRecord对应一个docker container
type LBController interface {
	StartLB(model.LBPolicy) error
	StopLB(model.LBPolicy) error
	DeleteLB(model.LBPolicy) error
	UpdateLB(model.LBPolicy) error
}

const (
	prefixDocker = "adc-haproxy-"
	imageName    = "hub.htres.cn/pub/haproxy:1.8"
)

//Controller lbagent manager
type Controller struct {
	WorkDir       string
	dockerClient  *dockercli.Client
	store         LBPolicyStore
	keepalivedCfg *KeepalivedConfig
}

// NewController create new Agent
func NewController(config *Config, store LBPolicyStore) (*Controller, error) {
	cli, err := dockercli.NewEnvClient()
	if err != nil {
		sysLogger.Errorf("Create haproxy controller failed, reason: %v", err)
		return nil, err
	}

	controller := &Controller{
		WorkDir:       config.WorkDir,
		dockerClient:  cli,
		store:         store,
		keepalivedCfg: config.KeepalivedCfg,
	}
	return controller, nil
}

// StartLB implement LBController
func (a *Controller) StartLB(policy model.LBPolicy) error {
	id := policy.GetID()
	needStart := false
	val, ok := a.store.GetByID(id)
	if !ok {
		a.store.Add(policy)
		needStart = true
	} else {
		if !val.Equals(policy) {
			needStart = true
			a.store.Add(policy)
		}
	}

	if needStart {
		// 如果开启高可用，先检查keepalived是否存在
		if a.keepalivedCfg != nil {
			if err := a.checkKeepalived(); err != nil {
				sysLogger.Errorf("check keepalived status failed, reason: %s", err)
				a.store.Delete(policy)
				return err
			}
		}

		containerID, err := a.doAddLB(policy)
		if err != nil {
			// 如果启动失败，删除store中的存储
			a.store.Delete(policy)
			return err
		}
		sysLogger.Debugf("create container with id %s", containerID)
	}
	return nil
}

// DeleteLB implement LBController
func (a *Controller) DeleteLB(policy model.LBPolicy) error {
	a.store.Delete(policy)
	return a.DelProxy(policy.GetID())
}

// UpdateLB implement LBController
func (a *Controller) UpdateLB(policy model.LBPolicy) error {
	return nil
}

// StopLB implement LBController
func (a *Controller) StopLB(policy model.LBPolicy) error {
	a.store.Delete(policy)
	return a.StopProxyByName(policy.GetID())
}

// doAddLB 根据policy生成container需要的配置文件
// 启动Haproxy container，实现负载均衡的添加
func (a *Controller) doAddLB(policy model.LBPolicy) (string, error) {
	filename := fmt.Sprintf("%s-haproxy.cfg", policy.GetID())
	configFilePath := filepath.Join(a.WorkDir, filename)

	// todo: 如果文件夹不存在则先创建文件夹
	file, err := os.Create(configFilePath)
	if err != nil {
		sysLogger.Errorf("Create haproxy config file failed, reason: %v", err)
		return "", err
	}

	if err := WriteHaproxyCfg(file, policy); err != nil {
		sysLogger.Errorf("Write haproxy config file failed, reason: %v", err)
		return "", err
	}
	if err := file.Close(); err != nil {
		sysLogger.Errorf("Close haproxy config file failed, reason: %v", err)
		return "", err
	}

	ports := []int32{}
	if policy.Record.Type == model.TypeIP {
		ports = append(ports, policy.Record.Port)
	} else {
		//如果是域名,则代理80和443端口
		if !portsContains(ports, 80) {
			ports = append(ports, 80)
		}
		if !portsContains(ports, 443) {
			ports = append(ports, 443)
		}
	}

	cid, err := a.RunProxy(configFilePath, policy, true)
	if err != nil {
		sysLogger.Errorf("Failed Start Proxy, reason: %v", err)
		return "", err
	}
	sysLogger.Infof("Successful add new loadbalancer %s", policy.GetID())
	return cid, err
}

// checkKeepalived 检查keepalived 容器是否存在，如果不存在就启动容器
func (a *Controller) checkKeepalived() error {
	id, err := getKeepalivedContainer(a.dockerClient)
	if err != nil {
		return err
	}
	if len(id) != 0 {
		return nil
	}

	cfgPath, err := createKeepalivedConfigFile(a.WorkDir, a.keepalivedCfg)
	if err != nil {
		return err
	}

	_, err = startKeepaliveContainer(a.dockerClient, cfgPath)
	if err != nil {
		return err
	}
	return nil
}

func portsContains(s []int32, t int32) bool {
	for _, a := range s {
		if a == t {
			return true
		}
	}
	return false
}
