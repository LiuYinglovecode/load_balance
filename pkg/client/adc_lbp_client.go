package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"code.htres.cn/casicloud/alb/pkg/model"
)

const (
	LBP_URL_SUFFIX = "/1/lbp"
)

type AdcLBPClient struct {
	lbpUrl string
	client *http.Client
}

func NewAdcLBPClient(url string, client *http.Client) *AdcLBPClient {
	if client == nil {
		client = http.DefaultClient
	}

	return &AdcLBPClient{lbpUrl: url, client: client}

}

type AdcLBPClientInterface interface {
	QuerySingleLBPById(id int64) string
	QueryAllLBP(ip string, DomainRegex string) string
	CreateLBP(LBPModel *model.LBPool) string
	UpdateLBP(id int64, param string) string
	DeleteLBP(LBPId int64) string
}

//查询单个LBP
func (client *AdcLBPClient) QuerySingleLBPById(id int64) string {
	url := client.lbpUrl + LBP_URL_SUFFIX + "/" + strconv.FormatInt(id, 10)
	return buildRequst(http.MethodGet, url, nil, client.client)

}

//所有列表LBP
func (client *AdcLBPClient) QueryAllLBP(ip string, type_pool string) string {
	url := client.lbpUrl + LBP_URL_SUFFIX

	var param = ""
	if ip != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "ip=" + ip
	}
	if type_pool != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "type=" + type_pool
	}

	return buildRequst(http.MethodGet, url+param, nil, client.client)
}

//创建LBP
func (client *AdcLBPClient) CreateLBP(LBPModel *model.LBPool) string {

	modelByte, err := json.Marshal(LBPModel)
	if err != nil {
		result := model.NewApiResult(model.ErrFormat, "参数解析异常", nil)
		resultByte, _ := json.Marshal(result)
		return string(resultByte)
	}

	fmt.Println("参数：", string(modelByte))
	return buildRequst(http.MethodPost, client.lbpUrl+LBP_URL_SUFFIX, bytes.NewBuffer(modelByte), client.client)

}

//修改LBP
func (client *AdcLBPClient) UpdateLBP(id int64, param string) string {
	if param == "" {
		result := model.NewApiResult(model.ErrParam, "修改参数不能为空", nil)
		resultByte, _ := json.Marshal(result)
		return string(resultByte)
	}
	url := client.lbpUrl + LBP_URL_SUFFIX + "/" + strconv.FormatInt(id, 10)
	return buildRequst(http.MethodPut, url, bytes.NewBuffer([]byte(param)), client.client)

}

//删除LBP
func (client *AdcLBPClient) DeleteLBP(LBPId int64) string {
	url := client.lbpUrl + LBP_URL_SUFFIX + "/" + strconv.FormatInt(LBPId, 10)
	return buildRequst(http.MethodDelete, url, nil, client.client)

}
