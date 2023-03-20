package query

import (
	"context"
	"time"

	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (op *AdminDB) ListTourPackages() (*model.ListResult, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	filter := bson.M{}

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
		Total: len(tourList),
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
		Total: len(tourList),
	}
	return data, nil
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
