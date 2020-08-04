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
	Openid    string  `json:"openid"`
	NickName  string  `json:"nick_name"`
}

func (Implement) SavePerson(ctx *gin.Context) {
	person := new(Person)

	if err := ctx.ShouldBindWith(person, binding.JSON); err != nil {
		log.Printf("%+v", xerrors.Errorf("bind json failed: %w", err))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
		})
		return
	}

	_, err := dbb.DB.SavePerson(person.Openid, person.NickName)
	if err != nil {
		log.Printf("%+v", xerrors.Errorf("db save person failed: %w", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
