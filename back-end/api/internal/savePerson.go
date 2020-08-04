package internal

import "github.com/gin-gonic/gin"

type Person struct {
	Openid    string  `json:"openid"`
	NickName  string  `json:"nick_name"`
}

func (Implement) SavePerson(ctx *gin.Context) {
}
