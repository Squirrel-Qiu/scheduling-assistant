package internal

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"os"
)

func (Implement) Download(ctx *gin.Context) {
	rotaId := ctx.Param("rotaId")

	filePath := "..." + rotaId + ".xlsx"

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Printf("%+v", xerrors.Errorf("file is not exist: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "值班表文件不存在",
		})
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Set("Content-Disposition", `attachment; filename="值班表.xlsx"`)
	ctx.Header("Content-Type", "application/octet-stream; charset=utf-8")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File(filePath)
}
