package internal

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
)

func (Implement) GetRotas(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	rotas, err := dbb.DB.GetRotas(openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db get rotas failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 1,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"rotas": rotas,
	})
}
