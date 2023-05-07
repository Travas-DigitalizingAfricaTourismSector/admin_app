package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"travas_admin/model"

	"travas_admin/internal/config"
	"travas_admin/internal/encrypt"
	"travas_admin/internal/query"
	"travas_admin/internal/token"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	App *config.Tools
	DB  query.Repo
}

func NewAdmin(app *config.Tools, db *mongo.Client) *Admin {
	return &Admin{
		App: app,
		DB:  query.NewAdminDB(app, db),
	}
}

// ProcessLogin : this method will help to parse, verify, and as well as authenticate the user
// login details, and it also helps to generate jwt token for restricted routers

func (op *Admin) ListOperators() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		nameParams := ctx.Query("name")
		page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
		limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)

		if page == 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 50
		}

		list, err := op.DB.ListOperators(page, limit, nameParams)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"data": list})
	}
}

func (op *Admin) ProcessLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}
		email := ctx.Request.Form.Get("email")
		password := ctx.Request.Form.Get("password")

		regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		ok := regMail.MatchString(email)

		if ok {
			res, checkErr := op.DB.VerifyUser(email)
			if checkErr != nil {
				_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unregistered user %v", checkErr))
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unregistered user"})
				return
			}

			id := (res["_id"]).(primitive.ObjectID)
			inputPass := (res["password"]).(string)
			// compName := (res["company_name"]).(string)

			verified, err := encrypt.Verify(password, inputPass)
			if err != nil {
				_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("cannot verify user details"))
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details"})
				return
			}
			switch {
			case verified:
				cookieData := sessions.Default(ctx)

				userInfo := model.UserInfo{
					ID:       id,
					Email:    email,
					Password: password,
					// CompanyName: compName,
				}
				cookieData.Set("info", userInfo)

				if err := cookieData.Save(); err != nil {
					log.Println("error from the session storage")
					_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
					return
				}
				// generate the jwt token
				// t1, t2, err := token.Generate(email, id)
				t1, _, err := token.Generate(email, id)

				if err != nil {
					_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token no generated : %v ", err))
				}

				cookieData.Set("token", t1)

				if err := cookieData.Save(); err != nil {
					log.Println("error from the session storage")
					_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
					return
				}

				// var tk map[string]string
				// tk := map[string]string{"t1": t1, "t2": t2}

				// // update the database adding the token to user database
				// _, updateErr := op.DB.UpdateInfo(userInfo.ID, tk)
				// if updateErr != nil {
				// 	_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("unregistered user %v", updateErr))
				// 	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect login details"})
				// 	return
				// }

				ctx.JSON(http.StatusOK, gin.H{
					"message": "Welcome to user homepage",
					"email":   email,
					"id":      id,
					// "company_name":  compName,
					"session_token": t1,
				})
			case !verified:
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details"})
				return
			}

		}
	}
}

func (op *Admin) ProcessReviews() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		response, err := op.DB.GetReviews()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"data": response})
	}
}
