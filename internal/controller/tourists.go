package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (op *Admin) ListTourists() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		list, err := op.DB.FindAllTourists()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"Tourists": list})
	}
}

func (op *Admin) GetBookingsSum() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		sum, err := op.DB.SumAllBookings()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"number_of_bookings": sum})
	}
}

func (op *Admin) GetRequestedTourSum() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		sum, err := op.DB.SumAllRequestedTours()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"requested_tours_sum": sum})
	}
}
