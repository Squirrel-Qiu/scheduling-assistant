package internal

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (impl *Implement) GetFrees(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	// 解析参数请求
	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "rotaId错误",
		})
		return
	}

	frees, err := impl.DB.GetFrees(openid, rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db get frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 2,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"frees": frees,
	})
}
