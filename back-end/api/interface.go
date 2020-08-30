package api

import (
	"github.com/gin-gonic/gin"

	"schedule/api/internal"
	"schedule/dbb"
	"schedule/snowid"
	"schedule/wechatid"
)

type Api interface {
	Login(ctx *gin.Context)
	SavePerson(ctx *gin.Context)

	NewRota(ctx *gin.Context)
	GetRotas(ctx *gin.Context)
	GetJoin(ctx *gin.Context)

	ChooseFree(ctx *gin.Context)
	GetFrees(ctx *gin.Context)

	Generate(ctx *gin.Context)
	Download(ctx *gin.Context)

	DeleteRota(ctx *gin.Context)
}

func New(dbInstance dbb.DBApi, openid wechatid.OpenidGetter, rotaId snowid.RotaidGetter) Api {
	return &internal.Implement{DB: dbInstance, OpenidGetter: openid, RotaidGetter: rotaId}
}
