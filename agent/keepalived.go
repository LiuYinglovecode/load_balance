package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockercli "github.com/docker/docker/client"
)

const (
	keepaliveContainerName = "alb-keepalived"
	keepaliveImageName     = "hub.htres.cn/pub/keepalived:2.0.15"
)

// createKeepalivedConfigFile 创建Keepalived配置文件，
// 接受工作目录路径， 返回配置文件路径
func createKeepalivedConfigFile(workdir string, k *KeepalivedConfig) (string, error) {
	filename := fmt.Sprintf("%s.cfg", keepaliveContainerName)
	configFilePath := filepath.Join(workdir, filename)

	// docker volume mount 需要使用绝对路径
	if !filepath.IsAbs(configFilePath) {
		absPath, err := filepath.Abs(configFilePath)
		if err != nil {
			return "", err
		}
		configFilePath = absPath
	}

	// todo: 如果文件夹不存在则先创建文件夹
	file, err := os.Create(configFilePath)
	if err != nil {
		sysLogger.Errorf("Create keepalived config file failed, reason: %v", err)
		return "", err
	}

	if err := WriteKeepalivedCfg(file, *k); err != nil {
		sysLogger.Errorf("Write keepalived config file failed, reason: %v", err)
		return "", err
	}
	if err := file.Close(); err != nil {
		sysLogger.Errorf("Close keepalived config file failed, reason: %v", err)
		return "", err
	}
	return configFilePath, nil
}

// startKeepaliveContainer 启动 keepalived 容器
func startKeepaliveContainer(cli *dockercli.Client, configFilePath string) (string, error) {
	ctx := context.Background()

	out, err := cli.ImagePull(ctx, keepaliveImageName, types.ImagePullOptions{})
	if err != nil {
		sysLogger.Errorf("Image pull failed, image_name: %s, reason: %v", keepaliveImageName, err)
		return "", err
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: keepaliveImageName,
		Tty:   true,
		// see https://github.com/osixia/docker-keepalived#fix-docker-mounted-file-problems
		Cmd: []string{"--copy-service"},
	}, &container.HostConfig{
		Binds:       []string{configFilePath + ":/container/service/keepalived/assets/keepalived.conf"},
		AutoRemove:  false,
		CapAdd:      []string{"NET_ADMIN"},
		NetworkMode: "host",
	}, nil, keepaliveContainerName)

	if err != nil {
		sysLogger.Errorf("Container create failed, container_name: %s, reason: %v", keepaliveContainerName, err)
		return "", err
	}

	sysLogger.Infof("Successful create keepalived container, container_id: %s", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		sysLogger.Errorf("Container start failed, container_name: %s, reason: %v", keepaliveContainerName, err)
		return "", err
	}

	sysLogger.Infof("Successful start keepalived container, container_id: %s", resp.ID)

	return resp.ID, nil
}

// getKeepalivedContainer 如果keepalived容器已经启动，返回容器id
// 否则返回空字符串
func getKeepalivedContainer(cli *dockercli.Client) (string, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Contains(n, keepaliveContainerName) {
				return c.ID, nil
			}
		}
	}
	return "", nil
}
