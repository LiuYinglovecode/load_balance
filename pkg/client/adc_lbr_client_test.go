package client

import (
	//"encoding/json"
	"testing"
	//"code.htres.cn/casicloud/alb/pkg/model"
)

const (
	LBR_URL = "http://127.0.0.1:8080"
)

//func TestAdcLbrClient_QuerySingleLBR(t *testing.T) {
//	cli := NewAdcLBRClient(LBR_URL, nil)
//	result := cli.QueryLBRById(0)
//	t.Log("result:" + result)
//}

func TestAdcLbrClient_QueryAllLBR(t *testing.T) {
	cli := NewAdcLBRClient(LBR_URL, nil)
	result := cli.QueryAllLBR("", "", "", "")
	t.Log("result:" + result)
}

//func TestAdcLbrClient_DeleteLbr(t *testing.T) {
//	cli := NewAdcLBRClient(LBR_URL, nil)
//	result := cli.DeleteLBR(6)
//	t.Log("result:" + result)
//}

//func TestAdcLbrClient_UpdateLBR(t *testing.T) {
//	cli := NewAdcLBRClient(LBR_URL, nil)

//	param := map[string]string{"ip": "127.0.0.1", "name": "描述测试"}

//	b, err := json.Marshal(param)
//	if err != nil {
//		return
//	}
//	result := cli.UpdateLBR(8, string(b))
//	t.Log("result:" + result)

//}

//func TestAdcLbrClient_CreateLBR(t *testing.T) {
//	cli := NewAdcLBRClient(LBR_URL, nil)
//	request := model.NewLBRecordIP(1223333, "127.0.0.1", 3301)
//	result := cli.CreateLBR(&request)
//	t.Log("result:" + result)

//}
