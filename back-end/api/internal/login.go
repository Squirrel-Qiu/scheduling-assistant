package internal

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"

	"schedule/middleware"
	"schedule/wechatid"
)

func (impl *Implement) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	code := new(wechatid.Code)

	if err := ctx.ShouldBindWith(code, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "请求参数错误",
		})
		return
	}

	openid, err := impl.OpenidGetter.GetOpenId(code)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("get openid failed: %w", err))
		// errcode:-1(系统繁忙); 40029(invalid code); 45011(频率限制)
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg":    "用户认证失败",
		})
		return
	}

	_, err = impl.DB.Login(openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db login failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	session.Set("openid", openid)

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
