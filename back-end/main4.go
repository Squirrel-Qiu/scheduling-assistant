package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/openid", getOpenID)
	http.ListenAndServe(":2020", nil)
}

func getOpenID(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		return
	}

	var codeMap map[string]string
	err := json.NewDecoder(request.Body).Decode(&codeMap)
	if err != nil {
		return
	}
	defer request.Body.Close()

	code := codeMap["code"]
	openid, err := sendWxAuthAPI(code)
	if err != nil {
		return
	}
	fmt.Println("my openid", openid)
}

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code"
	appID           = "你的AppID"
	appSecret       = "你的AppSecret"
)

func sendWxAuthAPI(code string) (string, error) {
	url := fmt.Sprintf(code2sessionURL, appID, appSecret, code)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	var wxMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&wxMap)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return wxMap["openid"], nil
}


