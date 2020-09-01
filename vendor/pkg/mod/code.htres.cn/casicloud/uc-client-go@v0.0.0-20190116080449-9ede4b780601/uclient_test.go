package uclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uclient_New(t *testing.T) {
	client := http.DefaultClient

	config := &UConfig{
		BaseURL: "http://auth.mst.casicloud.com/1/user",
	}

	uclient := NewUClient(config, client, nil)
	assert.Equal(t, uclient.config, config, "config not same")
}

func Test_uclient_Login(t *testing.T) {
	credential := UserCredential{
		UserName: "16608078741",
		Password: "123qwe",
	}

	client := http.DefaultClient

	config := &UConfig{
		BaseURL:      "http://auth.mst.casicloud.com/1",
		ClientID:     "fl4leakjuwvy8qp4",
		ClientSecret: "f32658f72c304e3885025f9de863afad",
	}

	uclient := NewUClient(config, client, nil)

	authInfo, err := uclient.Login(credential)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, authInfo)
	assert.Equal(t, authInfo.ClientID, config.ClientID)

}

func Test_uclient_GetUserInfo(t *testing.T) {
	credential := UserCredential{
		UserName: "16608078741",
		Password: "123qwe",
	}

	client := http.DefaultClient

	config := &UConfig{
		BaseURL:      "http://auth.mst.casicloud.com/1",
		ClientID:     "fl4leakjuwvy8qp4",
		ClientSecret: "f32658f72c304e3885025f9de863afad",
	}

	uclient := NewUClient(config, client, nil)

	authInfo, err := uclient.Login(credential)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, authInfo)
	assert.Equal(t, authInfo.ClientID, config.ClientID)

	userInfo, err := uclient.GetUserInfo(authInfo)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, userInfo)
	assert.Equal(t, userInfo.Mobile, credential.UserName)
}

func Test_uclient_GetOrgInfo(t *testing.T) {
	credential := UserCredential{
		UserName: "16608078741",
		Password: "123qwe",
	}

	client := http.DefaultClient

	config := &UConfig{
		BaseURL:      "http://auth.mst.casicloud.com/1",
		ClientID:     "fl4leakjuwvy8qp4",
		ClientSecret: "f32658f72c304e3885025f9de863afad",
	}

	uclient := NewUClient(config, client, nil)

	authInfo, err := uclient.Login(credential)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, authInfo)
	assert.Equal(t, authInfo.ClientID, config.ClientID)

	orgInfo, err := uclient.GetOrgInfo(10100001, authInfo)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, orgInfo)
}

func Test_uclient_RefreshToken(t *testing.T) {
	credential := UserCredential{
		UserName: "16608078741",
		Password: "123qwe",
	}

	client := http.DefaultClient

	config := &UConfig{
		BaseURL:      "http://auth.mst.casicloud.com/1",
		ClientID:     "fl4leakjuwvy8qp4",
		ClientSecret: "f32658f72c304e3885025f9de863afad",
	}

	uclient := NewUClient(config, client, nil)

	authInfo, err := uclient.Login(credential)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, authInfo)
	assert.Equal(t, authInfo.ClientID, config.ClientID)

	accessToken := authInfo.AccessToken
	expiresIn := authInfo.ExpiresIn

	err = uclient.RefreshToken(authInfo)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, authInfo)
	assert.Equal(t, authInfo.AccessToken, accessToken)
	assert.Equal(t, true, authInfo.ExpiresIn >= expiresIn)
}
