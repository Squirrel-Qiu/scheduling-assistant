package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"log"
	"net/http"
)

var SessionOption = sessions.Options{
	Path: "/",
	MaxAge: 5 * 60, // 5mins
}

func SessionChecker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		s := session.Get("openid")
		if s == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusForbidden,
			})
			ctx.Abort()
			return
		}
		openid := s.(string)

		session.Set("openid", openid)

		session.Options(SessionOption)

		if err := session.Save(); err != nil {
			log.Printf("%+v", xerrors.Errorf("save session failed: %w", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
			})
			ctx.Abort()
			return
		}

		ctx.Set("openid", openid)

		ctx.Next()
	}
}
