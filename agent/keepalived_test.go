package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	dockercli "github.com/docker/docker/client"
)

func TestKeepaliveContainerStartAndGet(t *testing.T) {
	cli, err := dockercli.NewEnvClient()
	assert.NoError(t, err)

	workDir := "."
	cfg := KeepalivedConfig{
		INet:            "eth0",
		VirutalRouterID: "51",
		State:           "MASTER",
		Priority:        100,
		UnicastSrcIP:    "127.0.0.1",
		UnicastPeer: []string{
			"127.0.0.2",
			"127.0.0.3",
		},
		VirtualIP: "8.8.8.8",
	}

	path, err := createKeepalivedConfigFile(workDir, &cfg)
	assert.NoError(t, err)

	id, err := startKeepaliveContainer(cli, path)
	assert.NoError(t, err)
	assert.NotEqual(t, "", id)

	id1, err := getKeepalivedContainer(cli)
	assert.NoError(t, err)
	assert.NotEqual(t, "", id1)

	os.Remove(path)
}
