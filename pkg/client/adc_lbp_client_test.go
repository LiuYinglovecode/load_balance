package client

import (
	//"encoding/json"
	"testing"
	//"code.htres.cn/casicloud/alb/pkg/model"
)

const (
	LBP_URL = "http://127.0.0.1:8080"
)

// func TestCreateLbP(t *testing.T) {
// 	modelLBP := NewAdcLBPClient(LBP_URL, nil)
// 	newModelLBR := model.LBPool{IP: "127.0.0.1", StartPort: 0}
// 	result := modelLBP.CreateLBP(&newModelLBR)
// 	t.Log("result:" + result)

// }

func TestQueryAllLBP(t *testing.T) {
	modelLBP := NewAdcLBPClient(LBP_URL, nil)
	result := modelLBP.QueryAllLBP("", "")
	t.Log("result:" + result)

}

// func TestQuerySingleLBP(t *testing.T) {
// 	modelLBP := NewAdcLBPClient(LBP_URL, nil)
// 	result := modelLBP.QuerySingleLBPById(3)
// 	t.Log("result:" + result)
// }

// func TestUpdateLBP(t *testing.T) {
// 	modelLBP := NewAdcLBPClient(LBP_URL, nil)
// 	param := map[string]interface{}{"ip": "127.0.0.2"}
// 	paramByte, _ := json.Marshal(param)
// 	result := modelLBP.UpdateLBP(5, string(paramByte))
// 	t.Log("result:" + result)

// }

// func TestDeleteLBP(t *testing.T) {
// 	modelLBP := NewAdcLBPClient(LBP_URL, nil)
// 	result := modelLBP.DeleteLBP(5)
// 	t.Log("result:" + result)

// }
