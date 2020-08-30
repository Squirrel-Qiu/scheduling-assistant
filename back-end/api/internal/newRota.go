package internal

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"

	"schedule/model"
)

func (impl *Implement) NewRota(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	rota := new(model.Rota)
	if err := ctx.ShouldBindWith(rota, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "请求参数错误",
		})
		return
	}

	if !newRotaCheck(rota) {
		log.Println("json content is wrong")
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg":    "限选不得小于班次",
		})
		return
	}

	rota.RotaId = impl.RotaidGetter.GetRotaId()

	_, err := impl.DB.NewRota(*rota, openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("when db new rota: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
		})
		return
	}

	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status": http.StatusBadRequest,
	//	})
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  0,
		"rota_id": strconv.FormatInt(rota.RotaId, 10),
	})
}
