package agent

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"testing"
)

var policy = model.LBPolicy{
	Record: model.LBRecord{
		IP:   model.NewADCString("192.168.100.200"),
		Port: 80,
		Type: model.TypeIP},
	Endpoints: []model.RealServer{
		{Name: "sever1", IP: "106.74.100.99", Port: 80},
		{Name: "sever2", IP: "106.74.100.98", Port: 80},
		{Name: "sever3", IP: "106.74.100.97", Port: 80},
	},
}

func TestController(t *testing.T) {
	home, err := homedir.Dir()
	assert.NoError(t, err)
	config := Config{
		WorkDir: home,
		AgentID: 1,
	}
	store, err := NewLBPolicyStore(&config)
	if err != nil {
		t.Error(err)
		return
	}
	c, err := NewController(&config, store)
	assert.NoError(t, err)

	t.Run("start lb", func(t *testing.T) {
		err := c.StartLB(policy)
		assert.NoError(t, err)
	})

	t.Run("stop lb", func(t *testing.T) {
		err := c.StopLB(policy)
		assert.NoError(t, err)
	})

	t.Run("stop lb", func(t *testing.T) {
		err := c.DeleteLB(policy)
		assert.NoError(t, err)
	})
}
