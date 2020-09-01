package agent

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// GetProxyList 获取当前启动的全部HaProxy容器名称
func (a *Controller) GetProxyList() ([]string, error) {
	containers, err := a.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Contains(n, prefixDocker) {
				names = append(names, n)
			}
		}
	}
	sysLogger.Infof("Get proxy container name list, names: %v", names)

	return names, nil
}

// GetProxyIDList 获取当前启动的全部HaProxy容器ID
func (a *Controller) GetProxyIDList() ([]string, error) {
	containers, err := a.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	cids := []string{}
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Contains(n, prefixDocker) {
				cids = append(cids, c.ID)
			}
		}
	}
	sysLogger.Infof("Get proxy container id list, cids: %v", cids)
	return cids, nil
}

// StopProxyByName 停止指定名称的容器
func (a *Controller) StopProxyByName(name string) error {
	return a.doStopAndRemoveContainer(name, false)
}

// RunProxy 启动容器
func (a *Controller) RunProxy(configFilePath string, policy model.LBPolicy, rm bool) (string, error) {
	name := policy.GetID()

	if rm {
		// 如果创建container时候auto remove = ture
		// 那么这里调用stop方法, 这样不好的地方在于不好查看容器日志
		// 否则调用del方法
		err := a.StopProxyByName(prefixDocker + name)
		if err != nil {
			sysLogger.Errorf("Stop container by name failed, container_name: %s, reason: %s", prefixDocker+name, err)
			return "", err
		}
	}

	ctx := context.Background()

	out, err := a.dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		sysLogger.Errorf("Image pull failed, container_name: %s, reason: %v", prefixDocker+name, err)
		return "", err
	}
	io.Copy(os.Stdout, out)

	resp, err := a.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   true,
		//Cmd:   []string{"sh", "-c", "echo hi, there! && sleep 3600"},
	}, &container.HostConfig{
		Binds:       []string{configFilePath + ":/usr/local/etc/haproxy/haproxy.cfg"},
		AutoRemove:  true,
		NetworkMode: "host",
	}, nil, prefixDocker+name)
	if err != nil {
		sysLogger.Errorf("Container create failed, container_name: %s, reason: %v", prefixDocker+name, err)
		return "", err
	}
	sysLogger.Infof("Successful create haproxy container, container_id: %s", resp.ID)

	if err := a.dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		sysLogger.Errorf("Container start failed, container_name: %s, reason: %v", prefixDocker+name, err)
		return "", err
	}
	sysLogger.Infof("Successful start haproxy container, container_id: %s", resp.ID)
	return resp.ID, nil
}

// DelProxy 删除容器
func (a *Controller) DelProxy(name string) error {
	return a.doStopAndRemoveContainer(name, true)
}

func (a *Controller) doStopAndRemoveContainer(name string, purge bool) error {
	containers, err := a.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		sysLogger.Errorf("List container failed, reason: %v", err)
		return err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.HasSuffix(n, name) {
				tm := 5 * time.Second
				if err := a.dockerClient.ContainerStop(context.TODO(), c.ID, &tm); err != nil {
					sysLogger.Errorf("Stop container by name failed, name: %s, reason: %v", n, err)
					return err
				}
				sysLogger.Infof("Stopped container, name: %v", n)

				if purge {
					if err := a.dockerClient.ContainerRemove(context.TODO(), c.ID, types.ContainerRemoveOptions{}); err != nil {
						sysLogger.Errorf("Stop and remove container by name failed, name: %s, reason: %v", n, err)
						return err
					}
					sysLogger.Infof("Deleted container, name: %v", n)
				}
				break
			}
		}
	}
	return err
}
