package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"
)

type Free struct {
	Frees  []int  `json:"frees,[]string"`
}

func (impl *Implement) ChooseFree(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	_, err := impl.DB.CheckNickName(openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("check nickName failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 5,
			"msg": "昵称存储失败,请点右上角重新进入小程序",
		})
		return
	}

	rotaId, err := strconv.ParseInt(ctx.Param("rotaId"), 10, 64)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("parse rotaId failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "rotaId错误",
		})
		return
	}

	frees := new(Free)

	if err := ctx.ShouldBindWith(frees, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "请求参数错误",
		})
		return
	}

	// check limitChoose
	limitChoose, err := impl.DB.GetLimitChoose(rotaId)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db query limit_choose failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "请求参数错误",
		})
		return
	}

	if len(frees.Frees) < limitChoose {
		log.Println("less than limitChoose")
		ctx.JSON(http.StatusOK, gin.H{
			"status": 3,
			"msg": fmt.Sprintf("请至少选择 %d 个时间段", limitChoose),
		})
		return
	}

	_, err = impl.DB.ChooseFree(openid, rotaId, frees.Frees)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db insert frees failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 4,
		})
		return
	}

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
