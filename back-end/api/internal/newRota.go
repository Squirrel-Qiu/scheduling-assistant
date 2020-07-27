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
			"status": http.StatusBadRequest,
		})
		return
	}

	if !newRotaCheck(rota) {
		log.Println("json content is wrong")
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	// UUID
	worker, err := snowFlakeByGo.NewWorker(0)
	if err != nil {
		log.Println("get uuid failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}
	rota.RotaId = worker.GetId()

	ok, err := dbb.DB.NewRota(*rota, openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("when db new rota: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
