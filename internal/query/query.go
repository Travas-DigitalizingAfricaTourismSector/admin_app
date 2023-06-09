package query

import (
	"context"
	"errors"
	"fmt"

	// "fmt"
	"regexp"
	"time"

	"travas_admin/internal/config"
	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminDB struct {
	App *config.Tools
	DB  *mongo.Client
}

func NewAdminDB(app *config.Tools, db *mongo.Client) Repo {
	return &AdminDB{
		App: app,
		DB:  db,
	}
}

// InsertUser : this will check for existing user and also insert new to the operator collection
func (op *AdminDB) InsertUser(user *model.Operator) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ok := regMail.MatchString(user.Email)
	if !ok {
		op.App.ErrorLogger.Println("invalid registered details")
		return false, 0, errors.New("invalid registered details")
	}

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := OperatorData(op.DB, "operators").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := OperatorData(op.DB, "operators").InsertOne(ctx, user)
			if insertErr != nil {
				op.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return true, 1, nil
		}
		op.App.ErrorLogger.Fatal(err)
	}
	return true, 2, nil
}

// VerifyUser : this method will verify the user login details store in the database
// and compare with the input details
func (op *AdminDB) VerifyUser(email string) (primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res bson.M

	filter := bson.D{{Key: "email", Value: email}}
	err := TouristsData(op.DB, "tourists").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			op.App.ErrorLogger.Println("no document found for this query")
			return nil, err
		}
		op.App.ErrorLogger.Fatalf("cannot execute the database query perfectly : %v ", err)
	}

	return res, nil
}

func (op *AdminDB) UpdateOperator(operatorID string, tk *model.Operator) (*model.Operator, error) {

	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// defer cancel()

	// filter := bson.D{{Key: "_id", Value: userID}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	// _, err := OperatorData(op.DB, "operators").UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	return false, err
	// }
	return nil, nil
}

// ValidTourRequest This query below is to get all the valid tour requested by the tourist
func (op *AdminDB) ValidTourRequest() ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{}}
	cursor, err := TouristsData(op.DB, "tourists").Find(ctx, filter)
	if err != nil {
		op.App.ErrorLogger.Fatal(err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			return
		}
	}(cursor, ctx)

	var tourist []bson.M
	if err = cursor.All(ctx, &tourist); err != nil {
		op.App.ErrorLogger.Fatal(err)
	}
	var validReq []primitive.M
	var rqTours []primitive.M
	for _, t := range tourist {
		rq := t["request_tour"].(primitive.A)
		for _, r := range rq {
			v := r.(primitive.M)
			date, _ := time.Parse("2006-01-02", v["start_time"].(string))
			currentDate, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
			if currentDate.Before(date) {
				newReq := map[string]interface{}{
					"first_name":   t["first_name"],
					"last_name":    t["last_name"],
					"request_tour": append(rqTours, v),
				}
				validReq = append(validReq, newReq)
			}
		}

	}

	return validReq, nil
}

func (op *AdminDB) ListOperators(page, limit int64, name string) (*model.ListResult, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := OperatorData(op.DB, "operators")

	var filter interface{}
	if name == "" {
		filter = bson.M{}

	} else {

		// regexPattern := fmt.Sprintf("^.*%s.*$", name)
		regexPattern := fmt.Sprintf("(?i).*%s.*", name)

		regexFilter := primitive.Regex{Pattern: regexPattern, Options: ""}

		filter = bson.M{"full_name": bson.M{"$regex": regexFilter}}

	}

	countChannel := make(chan int64)

	go func() {
		count, err := dataCollection.CountDocuments(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		countChannel <- count
	}()

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	var operatorList []*model.Operator

	if err = cur.All(context.TODO(), &operatorList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	data := &model.ListResult{
		Rows:  operatorList,
		Total: <-countChannel,
		Page:  page,
	}
	return data, nil
}

func (op *AdminDB) ListDashBoardOperators() ([]model.DashBoardOperator, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := OperatorData(op.DB, "operators")

	filter := bson.M{}

	cur, err := dataCollection.Find(ctx, filter)
	defer cur.Close(context.TODO())

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	var operatorList []model.DashBoardOperator

	if err = cur.All(context.TODO(), &operatorList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	return operatorList, nil
}

func (op *AdminDB) GetReviews() (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	// tours
	// tourGuide
	// operatorLiscence

	tourCollection := OperatorData(op.DB, "tours")
	tourChannel := make(chan int, 0)

	tourGuideCollection := OperatorData(op.DB, "tour_guide")
	tourGuideChannel := make(chan int, 0)

	operatorCollection := OperatorData(op.DB, "operators")
	operatorChannel := make(chan int, 0)

	// operatorCollection := OperatorData(op.DB, "operators")

	// filter := bson.M{}
	// tourFilter := bson.D{{Key: "isApproved", Value: false}}

	tourFilter := bson.M{"isApproved": false}
	go func() {

		cur, err := tourCollection.Find(ctx, tourFilter)

		var tourList []model.Tour

		if err = cur.All(context.TODO(), &tourList); err != nil {
			op.App.ErrorLogger.Fatal(err)
			// return nil, err
		}

		tourChannel <- len(tourList)
	}()

	go func() {

		cur, err := tourGuideCollection.Find(ctx, tourFilter)

		var tourList []model.Tour

		if err = cur.All(context.TODO(), &tourList); err != nil {
			op.App.ErrorLogger.Fatal(err)
			// return nil, err
		}

		tourGuideChannel <- len(tourList)

	}()
	go func() {

		cur, err := operatorCollection.Find(ctx, tourFilter)

		var operatorList []model.Operator

		if err = cur.All(context.TODO(), &operatorList); err != nil {
			op.App.ErrorLogger.Fatal(err)
			// return nil, err
		}

		operatorChannel <- len(operatorList)

	}()

	response := map[string]interface{}{
		"tours":          <-tourChannel,
		"toursGuides":    <-tourGuideChannel,
		"toursOperators": <-operatorChannel,
	}
	return response, nil
}

func (op *AdminDB) ApproveDeclineOperator(tg *model.Operator) (string, error) {

	dataCollection := OperatorData(op.DB, "operators")

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

func (op *AdminDB) GetOperator(operatorID string) (*model.Operator, error) {

	var operator *model.Operator
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	dataCollection := OperatorData(op.DB, "operators")

	operatorData, _ := primitive.ObjectIDFromHex(operatorID)
	filter := bson.M{"_id": operatorData}

	err := dataCollection.FindOne(ctx, filter).Decode(&operator)
	if err != nil {
		return nil, fmt.Errorf("error finding operator %v: %v", operatorID, err)
	}

	return operator, nil
}

func (op *AdminDB) ListReviewingOperators(requestData map[string]interface{}) (*model.ListResult, error) {

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
	dataCollection := OperatorData(op.DB, "operators")

	cur, err := dataCollection.Find(ctx, filter)

	if err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	defer cur.Close(context.TODO())

	var operatorList []*model.Operator

	if err = cur.All(context.TODO(), &operatorList); err != nil {
		op.App.ErrorLogger.Fatal(err)
		return nil, err
	}

	response := &model.ListResult{
		Rows:  operatorList,
		Total: int64(len(operatorList)),
	}
	return response, nil
}
