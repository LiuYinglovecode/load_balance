package uclient //import "code.htres.cn/casicloud/uc-client-go"

// UCClientInterface 用户中心访问接口
type UCClientInterface interface {
	// Login 用户登录
	Login(uc UserCredential) (authInfo *AuthInfo, err error)
	// 获取用户信息
	GetUserInfo(authInfo *AuthInfo) (userInfo *UserInfo, err error)
	// 根据企业ID获取企业信息
	GetOrgInfo(orgID int64, authInfo *AuthInfo) (orgInfo *OrgInfo, err error)
	// 刷新Token信息
	RefreshToken(authInfo *AuthInfo) error
}
