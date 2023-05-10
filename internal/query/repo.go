package query

import (
	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	ListOperators(page, limit int64, nameParams string) (*model.ListResult, error)
	UpdateOperator(operator_id string, operatorUpdate *model.Operator) (*model.Operator, error)
	ListDashBoardOperators() ([]model.DashBoardOperator, error)
	VerifyUser(email string) (primitive.M, error)
	ApproveDeclineOperator(guide *model.Operator) (string, error)
	GetOperator(ID string) (*model.Operator, error)
	ListReviewingOperators(map[string]interface{}) (*model.ListResult, error)

	// tours
	// InsertPackage(tour *model.Tour) (bool, error)
	// LoadTour(tourID primitive.ObjectID) (primitive.M, error)
	ListTourPackages(map[string]interface{}) (*model.ListResult, error)
	SumTourPackages() (int64, error)
	ListOperatorPackages(operatorID string) (*model.ListResult, error)
	ValidTourRequest() ([]primitive.M, error)
	ApproveDeclineTourPackage(tour *model.Tour) (string, error)
	GetTour(tourID string) (*model.Tour, error)

	// tourGuide
	InsertTourGuide(tg *model.TourGuide) (bool, error)
	FindTourGuide(operatorID primitive.ObjectID) ([]primitive.M, error)
	UpdateTourGuide(guideID string) error
	ListTourGuides(map[string]interface{}) (*model.ListResult, error)
	GetTourGuide(tourGuideID string) (*model.TourGuide, error)
	ListTourGuidesByOperator(operatorID string) (*model.ListResult, error)
	ApproveDeclineTourGuide(guide *model.TourGuide) (string, error)

	// tourists
	FindAllTourists(page, limit int64, nameParams string) (*model.ListResult, error)
	FindAllDashboardTourists() ([]model.DashBoardTourist, error)
	SumAllBookings() (*int, error)
	SumAllRequestedTours() (*int, error)

	// review
	GetReviews() (interface{}, error)
}
