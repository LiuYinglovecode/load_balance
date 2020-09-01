// Package uclient uc 客户端登录方式
package uclient //import "code.htres.cn/casicloud/uc-client-go"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type buildQueryFunc func() url.Values

// NewUClient new uc client
func NewUClient(config *UConfig, httpClient *http.Client, logger logrus.FieldLogger) *UClient {
	if logger == nil {
		logger = defaultLogger
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &UClient{
		client: httpClient,
		logger: logger,
		config: config,
	}
}

// Login interface implements of UCClientInterface Login method
func (u *UClient) Login(credential UserCredential) (authInfo *AuthInfo, err error) {
	if len(credential.UserName) == 0 || len(credential.Password) == 0 {
		// parameters nil
		return &AuthInfo{}, errors.New("invalid UserCredential")
	}

	jsonCredential, err := json.Marshal(credential)
	if err != nil {
		return &AuthInfo{}, errors.New("marshal credential failed")
	}
	// build url http://auth.mst.casicloud.com/1/user/auth?client_id=&sign=
	request, err := u.buildRequest("POST", "user/auth", func() url.Values {
		q := url.Values{}
		q.Add("client_id", u.config.ClientID)
		return q
	}, bytes.NewReader(jsonCredential))
	resp, err := u.client.Do(request) //发送请求
	if err != nil {
		return &AuthInfo{}, err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &AuthInfo{}, err
	}

	var authResult AuthResult
	err = json.Unmarshal(content, &authResult)
	if err != nil {
		return &AuthInfo{}, err
	}

	if authResult.Code == 200 {
		// ok return info
		return &(authResult.Info), nil
	}

	return &AuthInfo{}, fmt.Errorf("login failed with code %d with message %s",
		authResult.Code, authResult.Message)
}

// GetUserInfo implements UCClientInterface
func (u *UClient) GetUserInfo(authInfo *AuthInfo) (userInfo *UserInfo, err error) {
	request, err := u.buildRequest("GET", "user/get", func() url.Values {
		q := url.Values{}
		q.Add("client_id", u.config.ClientID)
		if authInfo != nil {
			q.Add("access_token", authInfo.AccessToken)
		}
		return q
	}, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(request) //发送请求
	if err != nil {
		return &UserInfo{}, err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &UserInfo{}, err
	}

	// src := string(content)
	// fmt.Println(src)
	var userInfoResult UserInfoResult
	err = json.Unmarshal(content, &userInfoResult)
	if err != nil {
		return &UserInfo{}, err
	}

	if userInfoResult.Code == 200 {
		// ok return info
		return &(userInfoResult.Info), nil
	}

	return &UserInfo{}, fmt.Errorf("get user info failed with code %d with message %s",
		userInfoResult.Code, userInfoResult.Message)
}

// GetOrgInfo implement UCClientInterface
func (u *UClient) GetOrgInfo(orgID int64, authInfo *AuthInfo) (orgInfo *OrgInfo, err error) {
	request, err := u.buildRequest("GET", "org/get", func() url.Values {
		q := url.Values{}
		q.Add("client_id", u.config.ClientID)
		if authInfo != nil {
			q.Add("access_token", authInfo.AccessToken)
		}
		return q
	}, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.client.Do(request) //发送请求
	if err != nil {
		return &OrgInfo{}, err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &OrgInfo{}, err
	}

	// src := string(content)
	// fmt.Println("org info ==================")
	// fmt.Println(src)
	// fmt.Println("===========================")
	var orgInfoResult OrgInfoResult
	err = json.Unmarshal(content, &orgInfoResult)
	if err != nil {
		return &OrgInfo{}, err
	}

	if orgInfoResult.Code == 200 {
		// ok return info
		return &(orgInfoResult.Info), nil
	}

	return &OrgInfo{}, fmt.Errorf("get org info failed with code %d with message %s",
		orgInfoResult.Code, orgInfoResult.Message)
}

// RefreshToken implements UCClientInterface
func (u *UClient) RefreshToken(authInfo *AuthInfo) error {
	request, err := u.buildRequest("GET", "token/refresh", func() url.Values {
		q := url.Values{}
		q.Add("client_id", u.config.ClientID)
		if authInfo != nil {
			q.Add("refresh_token", authInfo.RefreshToken)
			q.Add("grant_type", "refresh_token")
		}
		return q
	}, nil)
	if err != nil {
		return err
	}

	resp, err := u.client.Do(request) //发送请求
	if err != nil {
		return err
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// src := string(content)
	// fmt.Println("refresh token ==================")
	// fmt.Println(src)
	// fmt.Println("===========================")

	var authResult AuthResult
	err = json.Unmarshal(content, &authResult)
	if err != nil {
		return err
	}

	if authResult.Code == 200 {
		// ok return info
		authInfo.AccessToken = authResult.Info.AccessToken
		authInfo.RefreshToken = authResult.Info.RefreshToken
		authInfo.ClientID = authResult.Info.ClientID
		authInfo.ExpiresIn = authResult.Info.ExpiresIn
		authInfo.Scope = authResult.Info.Scope
		authInfo.UserID = authResult.Info.UserID
		return nil
	}

	return fmt.Errorf("login failed with code %d with message %s",
		authResult.Code, authResult.Message)

}

// buildRequest build rest request
func (u *UClient) buildRequest(method string, reqPath string, queryBuilder buildQueryFunc, body io.Reader) (*http.Request, error) {
	baseURL := fmt.Sprintf("%s/%s", u.config.BaseURL, reqPath)
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

	signed, err := signStringQuery(q.Encode(), u.config.ClientSecret)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = signed
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}
