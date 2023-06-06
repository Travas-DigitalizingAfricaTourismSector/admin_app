package main

import (
	"travas_admin/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, o controller.Admin) {
	r.MaxMultipartMemory = 10 << 20
	router := r.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.Default())

	cookieData := cookie.NewStore([]byte("travas"))
	router.Use(sessions.Sessions("session", cookieData))

	router.POST("/login", o.ProcessLogin())
	router.GET("/logout", o.ProcessLogOut())
	authRouter := r.Group("/admin")
	authRouter.Use(Authorization())

	{

		// TODO
		// GET Reviews
		// APPROVE TOURPACKAGE
		// APPROVE IDENTITY
		// APPROVE TOUR GUIDE

		// dasahboard lists
		authRouter.GET("/dashboard/operators", o.FindAllDashboardOperators())
		authRouter.GET("/dashboard/tourists", o.FindAllDashboardTourists())

		// operators
		authRouter.GET("/operators", o.ListOperators())
		// packages
		authRouter.GET("/packages", o.ListToursPackages())
		authRouter.GET("/packages/:operatorID", o.ListOperatorPackages())

		// guides
		authRouter.GET("/guides", o.ListTourGuides())
		authRouter.GET("/guides/:operatorID", o.ListOperatorTourGuides())

		// tour
		authRouter.GET("/tourists", o.ListTourists())

		authRouter.GET("/sum_packages", o.SumTourPackages())
		authRouter.GET("/sum_bookings", o.GetBookingsSum())
		authRouter.GET("/sum_requested_tours", o.GetRequestedTourSum())

		authRouter.GET("/reviews", o.ProcessReviews())
		authRouter.GET("/reviews/packages", o.ListToursPackagesToReview())
		authRouter.GET("/reviews/tour_guides", o.ListTourGuidesToReview())
		authRouter.GET("/reviews/operators", o.ListOperatorsToReview())
		authRouter.PATCH("/reviews/packages/:packageID", o.ApproveDeclineTourPackage())
		authRouter.PATCH("/reviews/tour_guides/:tourGuideID", o.ApproveDeclineTourGuide())
		authRouter.PATCH("/reviews/operators/:operatorID", o.ApproveDeclineOperator())

	}
}
