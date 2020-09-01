package client

import (
	"code.htres.cn/casicloud/alb/pkg/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 用于各模块间集成测试
func TestAlbClient_CreateLB(t *testing.T) {
	lbmcUrl := "http://127.0.0.1:8080"

	cli := NewAlbClient(lbmcUrl, nil)

	lbr := model.NewLBRecordIP("", "127.0.0.1", 18865)

	endpoints := []model.RealServer{
		{Name: "server1", IP: "172.17.60.113", Port:9008},
	}

	policy := model.LBPolicy{
		Record: lbr,
		Endpoints: endpoints,
	}

	req := model.NewLBRequest("12345", "test", model.ActionAdd, &policy)

	apiResult, err := cli.CreateLB(req)

	data := apiResult.Data.(map[string]interface{})

	fmt.Println(data["request_id"])
	assert.NoError(t, err)
	assert.NotEqual(t, "0", apiResult)
}

func TestAlbClient_QueryLB(t *testing.T) {
	lbmcUrl := "http://127.0.0.1:8080"

	cli := NewAlbClient(lbmcUrl, nil)
	apiResult, err := cli.QueryLB(148)
	assert.NoError(t, err)


	var status int32 = -1
	if data, ok := apiResult.Data.(map[string]interface{}); ok {
		if s, okk := data["status"]; okk {
			status = int32(s.(float64))
		}
	}

	assert.Equal(t, status, int32(2))
}