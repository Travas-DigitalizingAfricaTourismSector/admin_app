package query

import (
	"context"
	"fmt"
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

func (op *AdminDB) ListTourGuides(requestData map[string]interface{}) (*model.ListResult, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataValue, ok := requestData["isApproved"].(bool)
	filter := bson.M{}

	if !ok {

		filter = bson.M{}
	} else {
		filter = bson.M{
			"isApproved": dataValue,
		}
	}
	dataCollection := TourGuideData(op.DB, "tour_guide")

	cur, err := dataCollection.Find(ctx, filter)

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var tourGuideList []*model.TourGuide

	if err = cur.All(context.TODO(), &tourGuideList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	response := &model.ListResult{
		Rows:  tourGuideList,
		Total: int64(len(tourGuideList)),
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
		Total: int64(len(tourGuideList)),
	}
	return response, nil
}

func (op *AdminDB) ApproveDeclineTourGuide(tg *model.TourGuide) (string, error) {

	dataCollection := TourGuideData(op.DB, "tour_guide")

	filter := bson.M{"_id": tg.ID}

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
				"approvedBy": tg.ApprovedBy,
			},
		}
	}

	_, err := dataCollection.UpdateOne(context.TODO(), filter, updates)
	if err != nil {
		return "", err
	}

	return "successfully reviewed tourguide", nil
}

func (op *AdminDB) GetTourGuide(tourGuideID string) (*model.TourGuide, error) {

	var tourGuide *model.TourGuide
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := TourGuideData(op.DB, "tour_guide")

	filter := bson.M{"_id": tourGuideID}

	err := dataCollection.FindOne(ctx, filter).Decode(&tourGuide)
	if err != nil {
		return nil, fmt.Errorf("error finding tourGuide %v: %v", tourGuideID, err)
	}

	return tourGuide, nil
}
