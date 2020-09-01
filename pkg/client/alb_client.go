package client

import (
	"bytes"
	"code.htres.cn/casicloud/alb/pkg/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	lbRequestURLSuffix = "/1/lb"
)
// buildQueryFunc 生成请求参数
type buildQueryFunc func() url.Values

// AlbClient 封装LBMC的客户端操作
type AlbClient struct {
	LbmcURL string
	Client  *http.Client
}

// NewAlbClient 构造函数
func NewAlbClient(url string, client *http.Client) *AlbClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &AlbClient{
		LbmcURL: url,
		Client:  client,
	}
}

// CreateLB 向LBMC发送创建负载均衡请求, 返回 request_id
// 如果LBRequest中的action不是add, 会更改为add，然后发送
func (a *AlbClient) CreateLB(request *model.LBRequest) (*model.APIResult, error) {
	if request.Action != model.ActionAdd {
		request.Action = model.ActionAdd
	}
	return a.doSendRequest(request)
}

// StopLB 发送删除停止均衡请求
// note: 当前版本容器停止后自动删除
// 返回request id
func (a *AlbClient) StopLB(request *model.LBRequest) (*model.APIResult, error) {
	if request.Action != model.ActionStop{
		request.Action = model.ActionStop
	}
	return a.doSendRequest(request)
}

// QueryLB 查询请求处理状态
// todo: QueryLB 返回值应该加入IP, 给拥有资源池的用户(比如admin)使用
func (a *AlbClient) QueryLB(reqID int64) (*model.APIResult, error) {
	url :=  a.LbmcURL + lbRequestURLSuffix + "/" + strconv.FormatInt(reqID, 10)

	req, err := BuildRawRequest("GET", url,nil,nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var apiResult model.APIResult
	err = json.Unmarshal(body, &apiResult)
	if err != nil {
		return nil, err
	}
	return &apiResult, nil
}

func (a *AlbClient) doSendRequest(request *model.LBRequest) (*model.APIResult, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))

	req, err := BuildRawRequest("POST", a.LbmcURL + lbRequestURLSuffix, nil, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var apiResult model.APIResult
	err = json.Unmarshal(body, &apiResult)
	if err != nil {
		return nil, err
	}

	return &apiResult, nil
}

// BuildRawRequest 构建原始http请求
func BuildRawRequest(method string, baseURL string, queryBuilder buildQueryFunc, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, baseURL, body)
	if err != nil {
		return nil, err
	}

	var q url.Values
	if queryBuilder != nil {
		q = queryBuilder()
	} else {
		q = url.Values{}
	}

	request.URL.RawQuery = q.Encode()
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
