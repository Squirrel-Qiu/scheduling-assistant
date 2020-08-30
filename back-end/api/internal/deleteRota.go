package internal

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (impl *Implement) DeleteRota(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "rotaId错误",
		})
		return
	}

	_, err = impl.DB.DeleteRota(openid, rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db delete rota failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 2,
		})
		return
	}

	// 一般delete不会在err== nil的情况下返回false
	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status": http.StatusForbidden,
	//	})
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
	})
}
