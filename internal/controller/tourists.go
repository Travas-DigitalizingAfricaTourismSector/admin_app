package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (op *Admin) ListTourists() gin.HandlerFunc {
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
		list, err := op.DB.FindAllTourists(page, limit, nameParams)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"data": list})
	}
}

func (op *Admin) FindAllDashboardTourists() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Todo -> get all the tour request from the tourists collections
		//	and compare and check for the date with the present date

		list, err := op.DB.FindAllDashboardTourists()
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return

		}
		ctx.JSONP(http.StatusOK, gin.H{"tourists": list, "total": len(list)})
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
