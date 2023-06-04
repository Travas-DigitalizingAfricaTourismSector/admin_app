package controller

import (
	"errors"
	"fmt"
	"net/http"
	"travas_admin/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (op *Admin) ListToursPackages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		list, err := op.DB.ListTourPackages(nil)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourPackages": list})
	}
}

func (op *Admin) SumTourPackages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		list, err := op.DB.SumTourPackages()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourPackages_Sum": list})
	}
}

func (op *Admin) ListOperatorPackages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		operatorID := ctx.Param("operatorID")
		fmt.Println("GOT HERE", "operatorID", operatorID)
		list, err := op.DB.ListOperatorPackages(operatorID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourPackages": list})
	}
}

func (op *Admin) ListToursPackagesToReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		reviewFilter := map[string]interface{}{
			"isApproved":    false,
			"declineReason": "",
		}
		list, err := op.DB.ListTourPackages(reviewFilter)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourPackages": list})
	}
}

func (op *Admin) ApproveDeclineTourPackage() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		cookieData := sessions.Default(ctx)
		userInfo, ok := cookieData.Get("info").(model.UserInfo)

		if !ok {
			_ = ctx.AbortWithError(http.StatusNotFound, errors.New("cannot find admin id"))
		}

		packageID := ctx.Param("packageID")

		var input model.Tour
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error bindingJson %v", err)})
			return
		}

		if !input.IsApproved && input.DeclineReason == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "DeclineReason is required but empty."})
			return
		}

		tourPackage, err := op.DB.GetTour(packageID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error finding tour: %v", err)})
			return
		}

		if tourPackage.IsApproved {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "tour already approved"})
			return
		}

		data := &model.Tour{
			ID:            tourPackage.ID,
			IsApproved:    input.IsApproved,
			DeclineReason: input.DeclineReason,
			ApprovedBy:    userInfo.Email,
		}
		resp, err := op.DB.ApproveDeclineTourPackage(data)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"data": resp})
	}
}
