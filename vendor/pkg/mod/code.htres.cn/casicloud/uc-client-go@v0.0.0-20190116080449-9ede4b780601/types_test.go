package uclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// Test rest api client
func TestRest(t *testing.T) {
	// 测试http登录
	client := http.DefaultClient
	credential := &UserCredential{
		UserName: "16608078741",
		Password: "123qwe",
	}

	jsonCredential, err := json.Marshal(credential)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Printf("user credential is %s\n", string(jsonCredential))
	request, err := http.NewRequest("POST", "http://auth.mst.casicloud.com/1/user/auth?client_id=4g0ucoqrwtn92dxq&sign=ir834960bnjghze8343afajga", bytes.NewReader(jsonCredential))
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request) //发送请求
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	fmt.Println(string(content))

	var authResult AuthResult
	err = json.Unmarshal(content, &authResult)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(authResult.Info.AccessToken)

}

func TestMarshal(t *testing.T) {
	authInfoStr := `{"user_id":50000004260250,
	"access_token":"YB/5WnSYRGkgNYKdDScrTuQkUwt1JZ5WNyjvohVFVPwDelDSoNorf5WI0A57tYiqDhp2/wm2pSEM/MELxMtw5NS2MW2YuyfFqxs4m8zxuBo=",
	"expires_in":28800,
	"refresh_token":"bbe3f7fe42584a19b384a9b4e7afc2bf",
	"scope":"app_scope",
	"client_id":"4g0ucoqrwtn92dxq"}`

	var authInfo AuthInfo
	err := json.Unmarshal([]byte(authInfoStr), &authInfo)
	if err != nil {
		t.Error(err)
	}
	t.Logf("auth info is %s", authInfo.AccessToken)
}

func TestAuthAPIResult(t *testing.T) {
	authAPIReusltStr := `{"code":200,"msg":"","data":{"user_id":50000004260250,"access_token":"YB/5WnSYRGkgNYKdDScrTuQkUwt1JZ5WNyjvohVFVPwDelDSoNorf5WI0A57tYiqDhp2/wm2pSEM/MELxMtw5CxRfMpAvGErZSOfTUhBgGY=","expires_in":28800,"refresh_token":"1449742d901441c7acc4dc2911a394b6","scope":"app_scope","client_id":"4g0ucoqrwtn92dxq"}}`

	var authAPIResult AuthResult
	err := json.Unmarshal([]byte(authAPIReusltStr), &authAPIResult)
	if err != nil {
		t.Error(err)
	}
	t.Logf("auth info is %s", authAPIResult.Info.ClientID)
}
