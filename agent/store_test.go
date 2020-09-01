package agent

import (
	"testing"

	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/stretchr/testify/assert"
)

var store LBPolicyStore

func init() {
	var err error
	store, err = NewLBPolicyStore(&Config{
		WorkDir: "/tmp",
	})

	if err != nil {
		panic("init LBPolicyStore fail")
	}
}

func TestPolicyStoreSave(t *testing.T) {
	policies := []model.LBPolicy{
		model.LBPolicy{
			Record: model.LBRecord{
				ID:   123456,
				Type: model.TypeIP,
				IP:   model.NewADCString("192.168.100.10"),
				Port: 10240,
			},
			Endpoints: []model.RealServer{
				model.RealServer{
					Name: "web1",
					IP:   "192.168.100.120",
					Port: 3389,
				},
			},
		},
	}
	for _, v := range policies {
		err := store.Add(v)
		if err != nil {
			t.Error(err)
			return
		}
	}

	err := store.Save()
	if err != nil {
		t.Error(err)
		return
	}
	err = store.Load()
	if err != nil {
		t.Error(err)
		return
	}

	assert.True(t, store.Get()[0].Record.ID == policies[0].Record.ID)
	// fmt.Printf("%v\n", loadedPolicies)
}
