package query

import (
	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	InsertUser(user *model.Operator) (bool, int, error)
	VerifyUser(email string) (primitive.M, error)
	UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error)
	ListOperators() (*model.ListResult, error)
	ListDashBoardOperators() ([]model.DashBoardOperator, error)

	// tours
	// InsertPackage(tour *model.Tour) (bool, error)
	// LoadTour(tourID primitive.ObjectID) (primitive.M, error)
	ListTourPackages() (*model.ListResult, error)
	SumTourPackages() (int64, error)
	ListOperatorPackages(operatorID string) (*model.ListResult, error)

	ValidTourRequest() ([]primitive.M, error)

	// tourGuide
	InsertTourGuide(tg *model.TourGuide) (bool, error)
	FindTourGuide(operatorID primitive.ObjectID) ([]primitive.M, error)
	UpdateTourGuide(guideID string) error
	ListTourGuides() (*model.ListResult, error)
	ListTourGuidesByOperator(operatorID string) (*model.ListResult, error)

	// tourists
	FindAllTourists() ([]model.Tourist, error)
	FindAllDashboardTourists() ([]model.DashBoardTourist, error)
	SumAllBookings() (*int, error)
	SumAllRequestedTours() (*int, error)
}
