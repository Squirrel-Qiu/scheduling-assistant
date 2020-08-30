package internal

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (impl *Implement) GetJoin(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	joins, err := impl.DB.GetJoin(openid)
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
