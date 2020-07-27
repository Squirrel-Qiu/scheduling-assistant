package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
	"strconv"
)

func (Implement) ChooseFree(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	var frees []int
	if err := ctx.ShouldBindWith(&frees, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	// check limitChoose
	limitChoose, err := dbb.DB.GetLimitChoose(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query limit_choose failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}
	if len(frees) < limitChoose {
		log.Println("less than xxx")
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	ok, err := dbb.DB.ChooseFree(openid, rotaId, frees)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db insert frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
