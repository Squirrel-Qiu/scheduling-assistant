package internal

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
)

func (Implement) GetJoin(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	joins, err := dbb.DB.GetJoin(openid)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db get join failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 1,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
		"joins": joins,
	})
}
