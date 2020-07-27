package api

import (
	"github.com/gin-gonic/gin"
	"schedule/api/internal"
)

type Api interface {
	Login(ctx *gin.Context)

	NewRota(ctx *gin.Context)
	GetRotas(ctx *gin.Context)

	ChooseFree(ctx *gin.Context)
	GetFrees(ctx *gin.Context)
	Generate(ctx *gin.Context)

	DeleteRota(ctx *gin.Context)
}

func New() Api {
	return internal.Implement{}
}
