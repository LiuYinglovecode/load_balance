package agent

import (
	"context"
	"os"
	"testing"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
	"github.com/stretchr/testify/assert"
)

func TestEtcdInformer(t *testing.T) {
	startInProgressEtcdServer()

	config := &Config{
		AgentID:       1,
		KeepalivedCfg: &KeepalivedConfig{State: "MASTER"},
		Endpoints:     []string{"http://localhost:2379"},
	}

	config2 := &Config{
		AgentID:       1,
		KeepalivedCfg: &KeepalivedConfig{State: "BACKUP"},
		Endpoints:     []string{"http://localhost:2379"},
	}

	store, err := NewLBPolicyStore(config)
	assert.NoError(t, err)

	informer, err := NewInformer(config, store)
	assert.NoError(t, err)

	informer2, err := NewInformer(config2, store)
	assert.NoError(t, err)

	err = informer.Start()
	assert.NoError(t, err)

	err = informer.Stop()
	assert.NoError(t, err)

	err = informer.HeartBeat()
	assert.NoError(t, err)
	err = informer2.HeartBeat()
	assert.NoError(t, err)

	// 判断etcd中存在汇报数据
	etcdInformer := informer.(*etcdInformer)

	resp, err := etcdInformer.EtcdClient().Get(context.TODO(), InfomerKeyPrefx+"1", clientv3.WithPrefix())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Kvs))
}

const etcdDir = "default.etcd"

// 如果embed server启动失败，并且报如下错误:
// panic: codecgen version mismatch: current: 8, need 10.
// Re-generate file: ${GOPATH}/github.com/coreos/etcd/client/keys.generated.go
// 需要删除vendor目录对应的自动生成的文件, 通过以下命令实现：
// rm -f ${GOPATH}/code.htres.cn/casicloud/alb/vendor/github.com/coreos/etcd/client/keys.generated.go
// 详情: https://github.com/etcd-io/etcd/issues/8715
// 需要调用defer close() 方法关闭server
func startInProgressEtcdServer() *embed.Etcd {
	// start in progress etcd server
	deleteDataDirIfExists()
	cfg := embed.NewConfig()
	cfg.Dir = etcdDir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		panic(err)
	}
	return e
}

// 清除embed server 使用的数据文件
func deleteDataDirIfExists() {
	var err = os.RemoveAll(etcdDir)
	if err != nil {
		panic(err)
	}
}
