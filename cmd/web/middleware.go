package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"

	"travas_admin/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)

		fmt.Println("COOKIE DATA:", cookieData)
		// tokenString := cookieData.Get("token").(string)
		tokenString := cookieData.Get("token").(string)
		fmt.Println("tokenString:", tokenString)

		if tokenString == "" {
			_ = ctx.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}

		// fmt.Println(tokenString)

		parse, err := token.Parse(tokenString)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
		}
		ctx.Set("pass", tokenString)
		ctx.Set("id", parse.ID)
		ctx.Set("email", parse.Email)
		ctx.Next()
	}
}
