package agent

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Test_DockerRun(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	code, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("code is %d\n", code)
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

const haproxyCfg = `global
  daemon
  maxconn 2048

defaults
  mode tcp
  balance roundrobin
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms



listen listener-127.0.0.1:19999
  stick-table type ip size 200k expire 30m
  stick on src
  bind 127.0.0.1:19999
  server sever1 172.17.60.113:9008`

func TestController_RunProxy(t *testing.T) {
	name := "unittest"
	cfgFileName := "haproxy-demo.cfg"
	absPath, _ := filepath.Abs(cfgFileName)

	home, err := homedir.Dir()
	assert.NoError(t, err)
	config := Config{
		WorkDir: home,
		AgentID: 1,
	}
	c, err := NewController(&config, nil)
	assert.NoError(t, err)

	id, err := c.RunProxy(absPath, policy, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, id, "")

	names, err := c.GetProxyList()
	assert.NoError(t, err)
	assert.NotEmpty(t, names)

	ids, err := c.GetProxyIDList()
	assert.NoError(t, err)
	assert.NotEmpty(t, ids)

	err = c.DelProxy(name)
	assert.NoError(t, err)
}
