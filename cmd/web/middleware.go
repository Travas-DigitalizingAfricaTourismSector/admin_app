package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"

	"travas_admin/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)
		tokenString := cookieData.Get("token").(string)

		if tokenString == "" {
			_ = ctx.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}

		// fmt.Println(tokenString, "TOKENSTRING")
		// fmt.Println("TOKENSTRING")

		parse, err := token.Parse(tokenString)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
		}

		adminEmails := os.Getenv("ADMIN_EMAIL")

		emailList := strings.Split(adminEmails, ",")

		for i := 0; i < len(emailList); i++ {

			if parse.Email == emailList[i] {
				ctx.Set("pass", tokenString)
				ctx.Set("id", parse.ID)
				ctx.Set("email", parse.Email)
				ctx.Next()
				return
			}
		}
		_ = ctx.AbortWithError(http.StatusUnauthorized, gin.Error{Err: errors.New("you're unauthorized")})
		// return

	}
}
