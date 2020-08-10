package internal

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
	"schedule/middleware"
)

const URL = "https://api.weixin.qq.com/sns/jscode2session"

type Code struct {
	Appid  	  string  `json:"appid"`
	Secret 	  string  `json:"secret"`
	JsCode 	  string  `json:"js_code"`
	GrantType string  `json:"grant_type"`
}

type CommonError struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type LoginResponse struct {
	Openid     string  `json:"openid"`
	SessionKey string  `json:"session_key"`
	Unionid    string  `json:"unionid"`
}

type loginResponse struct {
	LoginResponse
	CommonError
}

func (Implement) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	code := new(Code)

	if err := ctx.ShouldBindWith(code, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "请求参数错误",
		})
		return
	}

	resp, err := getOpenid(*code)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("get openid failed: %w", err))
		// errcode:-1(系统繁忙); 40029(invalid code); 45011(频率限制)
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "用户认证失败",
		})
		return
	}

	_, err = dbb.DB.Login(resp.Openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db login failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	session.Set("openid", resp.Openid)

	session.Options(middleware.SessionOption)

	if err := session.Save(); err != nil {
		log.Printf("%+v", xerrors.Errorf("save session failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
	})
}

//"?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code"
func getOpenid(code Code) (LoginResponse, error) {
	request, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return LoginResponse{}, xerrors.Errorf("do request failed: %w", err)
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
		return LoginResponse{}, xerrors.Errorf("do client failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return LoginResponse{}, xerrors.New("WeChatServerError")
	}

	response := loginResponse{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	if err != nil {
		return LoginResponse{}, xerrors.Errorf("format response failed: %w", err)
	}

	if response.Errcode != 0 {
		return LoginResponse{}, xerrors.New(response.Errmsg)
	}

	return response.LoginResponse, nil
}
