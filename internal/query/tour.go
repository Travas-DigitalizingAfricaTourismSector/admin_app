package query

import (
	"context"
	"fmt"
	"time"

	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (op *AdminDB) ListTourPackages(reqData map[string]interface{}) (*model.ListResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataValue, ok := reqData["isApproved"].(bool)
	fmt.Println("dataValue::: ", dataValue, ok)
	filter := bson.M{}
	if !ok {

		filter = bson.M{}
	} else {
		filter = bson.M{
			"isApproved": dataValue,
		}
	}

	dataCollection := TourData(op.DB, "tours")

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}
	var tourList []*model.Tour

	if err = cur.All(context.TODO(), &tourList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	data := &model.ListResult{
		Rows:  tourList,
		Total: int64(len(tourList)),
	}
	return data, nil
}

func (op *AdminDB) SumTourPackages() (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	filter := bson.M{}

	dataCollection := TourData(op.DB, "tours")

	count, err := dataCollection.CountDocuments(ctx, filter)

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return 0, err
	}

	return count, nil
}

func (op *AdminDB) ListOperatorPackages(operatorID string) (*model.ListResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "operator_id", Value: operatorID}}

	dataCollection := TourData(op.DB, "tours")

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}
	var tourList []*model.Tour

	if err = cur.All(context.TODO(), &tourList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	data := &model.ListResult{
		Rows:  tourList,
		Total: int64(len(tourList)),
	}
	return data, nil
}

func (op *AdminDB) ApproveDeclineTourPackage(tg *model.Tour) (string, error) {

	var tour *model.Tour
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := TourData(op.DB, "tours")

	filter := bson.M{"_id": tg.ID}

	err := dataCollection.FindOne(ctx, filter).Decode(&tour)
	if err != nil {
		return "", err
	}
	var updates primitive.M

	if !tg.IsApproved {

		updates = bson.M{
			"$set": bson.M{
				"isApproved":    tg.IsApproved,
				"declineReason": tg.DeclineReason,
				"approvedBy":    tg.ApprovedBy,

			},
		}

	} else {
		updates = bson.M{
			"$set": bson.M{
				"isApproved": tg.IsApproved,
				"approvedBy":    tg.ApprovedBy,

			},
		}
	}

	_, err = dataCollection.UpdateOne(context.TODO(), filter, updates)
	if err != nil {
		return "", err
	}

	return "successfully reviewed tour", nil
}

func (op *AdminDB) GetTour(tourID string) (*model.Tour, error) {

	var tour *model.Tour
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := TourData(op.DB, "tours")

	filter := bson.M{"_id": tourID}

	err := dataCollection.FindOne(ctx, filter).Decode(&tour)
	if err != nil {
		return nil, fmt.Errorf("error finding tour %v: %v", tourID, err)
	}

	return tour, nil
}

type mongoPaginate struct {
	limit int64
	page  int64
}

func NewMongoPaginate(limit, page int64) *mongoPaginate {
	return &mongoPaginate{
		limit: (limit),
		page:  (page),
	}
}

func (mp *mongoPaginate) GetPaginatedOpts() *options.FindOptions {
	l := mp.limit
	skip := mp.page*mp.limit - mp.limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}
