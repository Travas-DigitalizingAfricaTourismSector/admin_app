package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (op *Admin) ListToursPackages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		list, err := op.DB.ListTourPackages()
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
		list, err := op.DB.ListOperatorPackages(operatorID)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"TourPackages": list})
	}
}
