package wechatid

import (
	"encoding/json"
	"net/http"

	"golang.org/x/xerrors"
)

const URL = "https://api.weixin.qq.com/sns/jscode2session"

type Code struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	JsCode    string `json:"js_code"`
	GrantType string `json:"grant_type"`
}

type LoginResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type OpenidGetter interface {
	GetOpenId(code *Code) (string, error)
}

type Wechat struct{}

//"?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code"
func (w Wechat) GetOpenId(code *Code) (string, error) {
	request, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return "", xerrors.Errorf("do request failed: %w", err)
	}

	query := request.URL.Query()

	query.Add("appid", code.Appid)
	query.Add("secret", code.Secret)
	query.Add("js_code", code.JsCode)
	query.Add("grant_type", code.GrantType)

	request.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()

	if err != nil {
		return "", xerrors.Errorf("do client failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", xerrors.New("WeChatServerError")
	}

	response := LoginResponse{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	if err != nil {
		return "", xerrors.Errorf("format response failed: %w", err)
	}

	if response.Errcode != 0 {
		return "", xerrors.New(response.Errmsg)
	}

	return response.Openid, nil
}
