package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/holdno/snowFlakeByGo"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
	"schedule/model"
)

func (Implement) NewRota(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	rota := new(model.Rota)
	if err := ctx.ShouldBindWith(rota, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "请求参数错误",
		})
		return
	}

	if !newRotaCheck(rota) {
		log.Println("json content is wrong")
		ctx.JSON(http.StatusOK, gin.H{
			"status": 2,
			"msg": "限选不得小于班次",
		})
		return
	}

	// rotaId
	worker, err := snowFlakeByGo.NewWorker(0)
	if err != nil {
		log.Println("get uuid failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 3,
			"msg": "实例化工作节点错误",
		})
		return
	}
	rota.RotaId = worker.GetId()

	_, err = dbb.DB.NewRota(*rota, openid)
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
		"status": 0,
		"rota_id": rota.RotaId,
	})
}
