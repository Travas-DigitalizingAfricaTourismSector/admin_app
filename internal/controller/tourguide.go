package controller

import (
	"errors"
	"fmt"
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
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find admin id"))
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
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find admin id"))
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
		cookieData := sessions.Default(ctx)
		fmt.Println("cookieData: ", cookieData)
		userInfo, ok := cookieData.Get("info").(model.UserInfo)

		if !ok {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find admin id"))
		}
		fmt.Println("userInfo, ok: ", userInfo, ok)

		list, err := op.DB.ListTourGuides(nil)
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

func (op *Admin) ListTourGuidesToReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		reviewFilter := map[string]interface{}{
			"isApproved":    false,
			"declineReason": "",
		}

		list, err := op.DB.ListTourGuides(reviewFilter)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourGuides": list})
	}
}

func (op *Admin) ApproveDeclineTourGuide() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)
		userInfo, ok := cookieData.Get("info").(model.UserInfo)

		if !ok {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find admin id"))
		}
		fmt.Println("userInfo, ok: ", userInfo, ok)

		tourID := ctx.Param("tourGuideID")

		var input model.TourGuide
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !input.IsApproved && input.DeclineReason == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "DeclineReason is required but empty."})
			return
		}

		guide, err := op.DB.GetTourGuide(tourID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error finding tourGuide: %v", err)})
			return
		}

		if guide.IsApproved {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "tourGuide already approved."})
			return
		}

		data := &model.TourGuide{
			ID:            tourID,
			IsApproved:    input.IsApproved,
			DeclineReason: input.DeclineReason,
			ApprovedBy:    userInfo.Email,
		}

		resp, err := op.DB.ApproveDeclineTourGuide(data)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"data": resp})
	}
}
