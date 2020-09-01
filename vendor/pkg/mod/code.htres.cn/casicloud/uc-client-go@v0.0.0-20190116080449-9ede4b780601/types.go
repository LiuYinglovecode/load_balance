// Package uclient 用户中心 app 客户端
// 根据《用户中心APP端接口文档v1.0.dox》生成相关类
package uclient //import "code.htres.cn/casicloud/uc-client-go"

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// APIResult api调用结果
type APIResult struct {
	Code    int    `json:"code"`
	Message string `json:"msg,omitempty"`
}

// UserCredential 用户登录凭证
type UserCredential struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// AuthInfo 登录后返回信息
type AuthInfo struct {
	AccessToken  string `json:"access_token"`  //分配的access token，有效期8小时
	ExpiresIn    int64  `json:"expires_in"`    //access token的有效期，单位为秒
	RefreshToken string `json:"refresh_token"` //分配的refresh token，用于刷新access token，有效期30天
	Scope        string `json:"scope"`         //授权作用域
	ClientID     string `json:"client_id"`     //分配的clientId
	UserID       int64  `json:"user_id"`       // 用户ID
}

// UserInfo 用户中心，用户相关信息
type UserInfo struct {
	UserID       int64  `json:"user_id"`       // 用户ID
	FullName     string `json:"fullname"`      // 用户姓名
	Account      string `json:"account"`       // 登录名
	ShortAccount string `json:"short_account"` // 短名
	Status       int16  `json:"status"`        //状态
	Email        string `json:"email"`         //邮箱
	Mobile       string `json:"mobile"`        //用户手机
	Phone        string `json:"phone"`         //电话
	Sex          string `json:"sex"`           //用户性别，0代表女，1代表男
	Picture      string `json:"picture"`       //用户头像
	OpenID       string `json:"open_id"`       // 用户OpenID
	OrgID        int64  `json:"org_id"`        //用户所属的企业ID
	OrgName      string `json:"org_name"`      //用户所属企业名称
	IsAdmin      bool   `json:"is_admin"`      //是否管理员， default false
	CreateTime   string `json:"create_time"`   // 创建时间，目前暂时用string处理
	UpdateTime   string `json:"update_time"`   //更新时间,目前暂时用string处理
	UpdateTimes  int32  `json:"update_times"`  //用户账号修改次数
}

// OrgInfo 企业信息
type OrgInfo struct {
	OrgID      int64  `json:"org_id"`            //企业ID
	Name       string `json:"name"`              //企业名称
	OpenID     string `json:"open_id"`           // 企业OpenID
	Province   string `json:"province"`          //企业所在省份
	City       string `json:"city"`              //企业所在城市
	County     string `json:"county"`            //企业所在区县
	Address    string `json:"address"`           //企业地址
	Industry   string `json:"industry"`          //企业所属一级行业
	Industry2  string `json:"industry2"`         //企业所属二级行业
	Contact    string `json:"connecter"`         //联系人 	FIXME: typo for json string connector
	Telephone  string `json:"tel"`               //联系人手机号
	Fax        string `json:"fax"`               //传真
	HomePhone  string `json:"homephone"`         //固定电话
	Email      string `json:"email"`             //邮箱
	Postcode   string `json:"postcode"`          //邮编
	ShortName  string `json:"abbreviation_name"` //企业简称
	CreateTime string `json:"create_time"`       // 创建时间，目前暂时用string处理
	UpdateTime string `json:"update_time"`       //更新时间,目前暂时用string处理
}

// AuthResult auth api result
type AuthResult struct {
	APIResult
	Info AuthInfo `json:"data"`
}

// UserInfoResult user info api result
type UserInfoResult struct {
	APIResult
	Info UserInfo `json:"data"`
}

// OrgInfoResult org info result
type OrgInfoResult struct {
	APIResult
	Info OrgInfo `json:"data"`
}

var defaultLogger = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: &logrus.TextFormatter{DisableColors: true},
	Level:     logrus.DebugLevel,
}

// UConfig uclient config
type UConfig struct {
	// REST URL base地址 http://auth.mst.casicloud.com/1
	BaseURL      string `json:"base_url" yaml:"base_url"`
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
}

// UClient uc app client implments UCClientInterface
type UClient struct {
	client *http.Client
	logger logrus.FieldLogger
	config *UConfig
}
