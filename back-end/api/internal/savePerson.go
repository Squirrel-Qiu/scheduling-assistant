package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"schedule/dbb"
)

type Person struct {
	NickName  string  `json:"nick_name"`
}

func (Implement) SavePerson(ctx *gin.Context) {
	openid := ctx.Value("openid").(string)

	person := new(Person)

	if err := ctx.ShouldBindWith(person, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg": "请求参数错误",
		})
		return
	}

	_, err := dbb.DB.SavePerson(openid, person.NickName)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db save person failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": 2,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": 0,
	})
}
