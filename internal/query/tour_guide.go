package query

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (op *AdminDB) InsertTourGuide(tg *model.TourGuide) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: tg.ID}}

	var res bson.M
	err := OperatorData(op.DB, "tour_guide").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, insertErr := OperatorData(op.DB, "tour_guide").InsertOne(ctx, &tg)
			if insertErr != nil {
				op.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return true, nil
		}
		op.App.ErrorLogger.Fatal(err)
	}
	return true, nil
}

func (op *AdminDB) FindTourGuide(operatorID primitive.ObjectID) ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "operator_id", Value: operatorID}}

	cursor, err := TourGuideData(op.DB, "tour_guide").Find(ctx, filter)
	if err != nil {
		op.App.ErrorLogger.Fatalf("error while searching for data : %v \n", err)

	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var res []bson.M
	if err := cursor.All(ctx, &res); err != nil {
		op.App.ErrorLogger.Fatal(err)
	}

	return res, nil
}

func (op *AdminDB) UpdateTourGuide(guideID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: guideID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "_id", Value: guideID}}}}
	opt := options.Update().SetUpsert(false)

	_, err := TourGuideData(op.DB, "tour_guide").UpdateOne(ctx, filter, update, opt)
	if err != nil {
		op.App.ErrorLogger.Fatal(err)
	}
	return nil
}

func (op *AdminDB) ListTourGuides() (*model.ListResult, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := TourGuideData(op.DB, "tour_guide")

	filter := bson.M{}

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	var tourGuideList []*model.TourGuide

	if err = cur.All(context.TODO(), &tourGuideList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	response := &model.ListResult{
		Rows:  tourGuideList,
		Total: len(tourGuideList),
	}
	return response, nil
}

func (op *AdminDB) ListTourGuidesByOperator(operatorID string) (*model.ListResult, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := TourGuideData(op.DB, "tour_guide")

	filter := bson.D{{Key: "operator_id", Value: operatorID}}

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	var tourGuideList []*model.TourGuide

	if err = cur.All(context.TODO(), &tourGuideList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	response := &model.ListResult{
		Rows:  tourGuideList,
		Total: len(tourGuideList),
	}
	return response, nil
}
