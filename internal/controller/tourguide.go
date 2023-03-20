package controller

import (
	"errors"
	"net/http"
	"strings"

	"travas_admin/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (op *Admin) AddTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if err := ctx.Request.ParseForm(); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}

		cookieData := sessions.Default(ctx)
		userInfo, ok := cookieData.Get("info").(model.UserInfo)

		if !ok {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find operator id"))
		}

		tourGuide := &model.TourGuide{
			OperatorID: userInfo.ID,
			ID:         primitive.NewObjectID().Hex(),
			Name:       ctx.Request.Form.Get("guide_name"),
			Bio:        ctx.Request.Form.Get("bio"),
		}

		ok, err := op.DB.InsertTourGuide(tourGuide)
		if !ok {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}
		ctx.JSONP(http.StatusOK, gin.H{"message": "New tour guide added"})
	}
}

func (op *Admin) GetTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)
		userInfo, ok := cookieData.Get("info").(model.UserInfo)

		if !ok {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find operator id"))
		}

		arrRes, err := op.DB.FindTourGuide(userInfo.ID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		// all the value to the id key should be added to the remove url
		ctx.JSONP(http.StatusOK, gin.H{"TourGuides": arrRes})
	}
}

func (op *Admin) DeleteTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		guideID := ctx.Param("id")
		guideID = strings.TrimSpace(guideID)
		err := op.DB.UpdateTourGuide(guideID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"message": "successfully remove tour guide"})

	}
}

func (op *Admin) SelectTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	Todo -> make enquire from the frontend dev before coding this up
		// find the right selected guide and add that to the tour packages as well

	}
}

func (op *Admin) ListTourGuides() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	Todo -> make enquire from the frontend dev before coding this up
		// find the right selected guide and add that to the tour packages as well

		list, err := op.DB.ListTourGuides()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourGuides": list})
	}
}

func (op *Admin) ListOperatorTourGuides() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	Todo -> make enquire from the frontend dev before coding this up
		// find the right selected guide and add that to the tour packages as well
		operatorID := ctx.Param("operatorID")
		list, err := op.DB.ListTourGuidesByOperator(operatorID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourGuides": list})
	}
}
