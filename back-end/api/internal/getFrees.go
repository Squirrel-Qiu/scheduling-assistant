package internal

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
	"strconv"
)

func (Implement) GetFrees(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	// 解析参数请求
	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	frees, err := dbb.DB.GetFrees(openid, rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db get frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"frees": frees,
	})
}