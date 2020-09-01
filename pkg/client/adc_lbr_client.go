package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"code.htres.cn/casicloud/alb/pkg/model"
)

const (
	//LBR
	LBRRequestURISuffix = "/1/lbr"
)

type AdcLBRClient struct {
	lbrUrl string
	Client *http.Client
}

func NewAdcLBRClient(url string, client *http.Client) *AdcLBRClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &AdcLBRClient{
		lbrUrl: url,
		Client: client,
	}
}

type AdcLBRClientInterface interface {
	CreateLBR(request *model.LBRecord) string
	UpdateLBR(lbrId int64, request *model.LBRecord) string
	DeleteLBR(lbrId int64) string
	QueryLBRById(lbrId int64) string
	QueryAllLBR(owner string, ip string, port string, domain string) string
}

//lbr相关的客户端请求
//创建lbr
func (a *AdcLBRClient) CreateLBR(request *model.LBRecord) string {

	b, err := json.Marshal(request)
	if err != nil {
		fmt.Println("参数解析异常:", err.Error())
		result := model.NewApiResult(model.ErrFormat, "参数解析异常", nil)
		resultStr, _ := json.Marshal(result)
		return string(resultStr)
	}
	fmt.Println("创建内容", string(b))

	url := a.lbrUrl + LBRRequestURISuffix
	return buildRequst(http.MethodPost, url, bytes.NewBuffer(b), a.Client)
}

//更新LBR
func (a *AdcLBRClient) UpdateLBR(lbrId int64, param string) string {
	fmt.Println("更新内容", param)
	url := a.lbrUrl + LBRRequestURISuffix + "/" + strconv.FormatInt(lbrId, 10)
	return buildRequst(http.MethodPut, url, bytes.NewBuffer([]byte(param)), a.Client)

}

//删除LBR
func (a *AdcLBRClient) DeleteLBR(lbrId int64) string {
	url := a.lbrUrl + LBRRequestURISuffix + "/" + strconv.FormatInt(lbrId, 10)
	return buildRequst(http.MethodDelete, url, nil, a.Client)

}

//查询单个LBR
func (a *AdcLBRClient) QueryLBRById(lbrId int64) string {
	url := a.lbrUrl + LBRRequestURISuffix + "/" + strconv.FormatInt(lbrId, 10)
	return buildRequst(http.MethodGet, url, nil, a.Client)

}

//查询所有的LBR
func (a *AdcLBRClient) QueryAllLBR(owner string, ip string, port string, domain string) string {
	url := a.lbrUrl + LBRRequestURISuffix

	var param = ""
	if owner != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "owner=" + owner
	}

	if ip != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "ip=" + ip
	}

	if port != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "port=" + port
	}
	if domain != "" {
		if param == "" {
			param += "?"
		} else {
			param += "&"
		}
		param += "domain=" + domain
	}

	return buildRequst(http.MethodGet, url+param, nil, a.Client)

}

func buildRequst(method string, urlInfo string, body io.Reader, modelClient *http.Client) string {
	fmt.Println("url:" + urlInfo + " method:" + method)

	req, err := http.NewRequest(method, urlInfo, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}

	resp, err := modelClient.Do(req)
	if err != nil {
		fmt.Println("请求异常:", err)
		result := model.NewApiResult(model.ErrNotFound, "接口异常", nil)
		resultByte, _ := json.Marshal(result)
		return string(resultByte)
	}

	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("响应内容：", string(s))
	return string(s)

}
