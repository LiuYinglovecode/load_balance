package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
	"github.com/jinzhu/gorm"
)

const etcdDir = "default.etcd"

var repoStub = LBAgentRepositoryStub{}

type LBAgentRepositoryStub struct {
}

func (*LBAgentRepositoryStub) Create(agent *model.LBAgent) error {
	panic("implement me")
}

func (*LBAgentRepositoryStub) GetByID(id int64) (*model.LBAgent, error) {
	panic("implement me")
}

func (*LBAgentRepositoryStub) Update(agent *model.LBAgent) error {
	panic("implement me")
}

func (*LBAgentRepositoryStub) DeleteByID(id int64) error {
	panic("implement me")
}

func (*LBAgentRepositoryStub) DB() *gorm.DB {
	panic("implement me")
}

func (*LBAgentRepositoryStub) List(conditions map[string]interface{}) ([]model.LBAgent, error) {
	return []model.LBAgent{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}, nil
}

func TestEtcdReadWriter_Enqueue(t *testing.T) {
	t.Run("test enqueue lbresponse into etcd", func(t *testing.T) {
		// start in progress etcd server
		server := startInProgressEtcdServer()
		defer server.Close()

		lbr := &model.LBRequest{User: model.NewADCString("12345")}
		wantKey := "1"
		wantValue, err := json.Marshal(lbr)

		ep := []string{"http://localhost:2379"}
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   ep,
			DialTimeout: 5 * time.Second})
		AssertNoError(t, err)

		e, err := NewEtcdMessageQueueHandler(ep, "", "", "", &repoStub)
		AssertNoError(t, err)

		err = e.Enqueue(wantKey, lbr)
		err = e.Enqueue(wantKey, lbr)
		AssertNoError(t, err)

		gr, err := client.Get(context.TODO(), common.LBReqKeyPrefix+wantKey, clientv3.WithPrefix())
		if len(gr.Kvs) == 0 {
			t.Errorf("want key: %s, value: %s, but got none", wantKey, wantValue)
		}

		i := 0
		for _, ev := range gr.Kvs {
			gotValue := ev.Value

			AssertKvEquals(t, string(gotValue), string(wantValue))
			i = i + 1
		}

		if i != 2 {
			t.Errorf("got %d, want 2", i)
		}
	})
}

func TestEtcdReadWriter_WatchAndDequeue(t *testing.T) {
	t.Run("test watch and dequeue", func(t *testing.T) {
		// start in progress etcd server
		server := startInProgressEtcdServer()
		defer server.Close()

		ep := []string{"http://localhost:2379"}
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   ep,
			DialTimeout: 5 * time.Second})
		AssertNoError(t, err)

		e, err := NewEtcdMessageQueueHandler(ep, "", "", "", &repoStub)
		AssertNoError(t, err)

		user := "12345"
		wantKey := common.LBReqKeyPrefix + "1"
		wantValue, err := json.Marshal(model.LBRequest{User: model.NewADCString(user)})
		AssertNoError(t, err)

		_, err = client.Put(context.TODO(), wantKey+"/123", string(wantValue))
		_, err = client.Put(context.TODO(), wantKey+"/234", string(wantValue))
		AssertNoError(t, err)

		req1, _ := e.WatchAndDequeue()
		if req1.User.String() != user {
			t.Errorf("got %s, want %s", req1.User.String(), user)
		}
		req2, _ := e.WatchAndDequeue()
		if req2.User.String() != user {
			t.Errorf("got %s, want %s", req2.User.String(), user)
		}
	})
}

func TestEtcdReadWriter_OnAgentChange(t *testing.T) {
	// start in progress etcd server
	server := startInProgressEtcdServer()
	defer server.Close()

	ep := []string{"http://localhost:2379"}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   ep,
		DialTimeout: 5 * time.Second})
	AssertNoError(t, err)

	qh, err := NewEtcdMessageQueueHandler(ep, "", "", "", &repoStub)
	AssertNoError(t, err)
	e := qh.(*EtcdReadWriter)
	go e.OnAgentChange()

	t.Run("test add agent", func(t *testing.T) {
		newQueueName := "i'm new"
		addEvent := &LBAgentEvent{LBAgentAdd, newQueueName, false}
		o := e.GetCurrentQueueNames()
		oldLen := len(o)

		c := GetLBAgentEventChannel()
		c <- addEvent

		// wait for onAgentChange()
		time.Sleep(1 * time.Second)

		n := e.GetCurrentQueueNames()
		newLen := len(n)

		if newLen-oldLen != 1 {
			t.Errorf("want 1, got %d", newLen-oldLen)
		}

		if !contains(n, common.LBReqKeyPrefix+newQueueName) {
			t.Errorf("want contains %s, but does not", newQueueName)
		}
	})

	t.Run("test delete agent", func(t *testing.T) {
		queueName := "1"
		deleteEvent := &LBAgentEvent{LBAgentDelete, queueName, false}
		o := e.GetCurrentQueueNames()
		oldLen := len(o)

		c := GetLBAgentEventChannel()
		c <- deleteEvent

		// wait for onAgentChange()
		time.Sleep(1 * time.Second)

		n := e.GetCurrentQueueNames()
		newLen := len(n)

		if oldLen-newLen != 1 {
			t.Errorf("want 1, got %d", oldLen-newLen)
		}

		if contains(n, queueName) {
			t.Errorf("do not want contains %s, but contains", queueName)
		}
	})

	t.Run("test delete with purge flag", func(t *testing.T) {
		queueName := "2"
		deleteEvent := &LBAgentEvent{LBAgentDelete, queueName, true}

		_ = e.Enqueue(queueName, &model.LBRequest{})
		_ = e.Enqueue(queueName, &model.LBRequest{})

		gr, err := client.Get(context.TODO(), common.LBReqKeyPrefix+queueName, clientv3.WithPrefix())
		if len(gr.Kvs) == 0 {
			t.Errorf("want 2, but got 0")
		}

		c := GetLBAgentEventChannel()
		c <- deleteEvent

		// wait for onAgentChange()
		time.Sleep(1 * time.Second)

		gr, err = client.Get(context.TODO(), queueName, clientv3.WithPrefix())
		if len(gr.Kvs) != 0 {
			fmt.Print(gr.Kvs)
			t.Errorf("want queue purged, but not")
		}
		AssertNoError(t, err)
	})
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("want no err, but got %v", err)
	}
}

func AssertKvEquals(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// 清除embed server 使用的数据文件
func deleteDataDirIfExists() {
	var err = os.RemoveAll(etcdDir)
	if err != nil {
		panic(err)
	}
}

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
