package query

import (
	"context"
	"fmt"
	"time"
	"travas_admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (td *AdminDB) FindAllTourists(page, limit int64, name string) (*model.ListResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	dataCollection := TouristsData(td.DB, "tourists")
	countChannel := make(chan int64)

	var filter interface{}
	if name == "" {
		filter = bson.M{}

	} else {
		// regexPattern := fmt.Sprintf("^.*%s.*$", name)
		regexPattern := fmt.Sprintf("(?i).*%s.*", name)

		regexFilter := primitive.Regex{Pattern: regexPattern, Options: ""}

		firstNameFilter := bson.M{"first_name": bson.M{"$regex": regexFilter}}
		lastNameFilter := bson.M{"last_name": bson.M{"$regex": regexFilter}}

		filter = bson.M{"$or": []bson.M{firstNameFilter, lastNameFilter}}

	}

	go func() {
		count, err := dataCollection.CountDocuments(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		countChannel <- count
	}()

	tourists := []model.Tourist{}
	cursor, err := dataCollection.Find(context.TODO(), filter, NewMongoPaginate(limit, page).GetPaginatedOpts())

	if err != nil {
		return nil, fmt.Errorf("cannot find document in the database %v ", err)
	}

	if err = cursor.All(ctx, &tourists); err != nil {
		return nil, err
	}

	data := &model.ListResult{
		Rows:  tourists,
		Total: <-countChannel,
		Page:  page,
	}
	return data, err
}

func (td *AdminDB) FindAllDashboardTourists() (tourists []model.DashBoardTourist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := TouristsData(td.DB, "tourists").Find(ctx, bson.D{{}})
	if err != nil {
		return tourists, fmt.Errorf("cannot find document in the database %v ", err)
	}

	if err = cursor.All(ctx, &tourists); err != nil {
		return nil, err
	}

	return tourists, err
}

func (td *AdminDB) SumAllBookings() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var tourists []model.Tourist

	bookingSum := 0

	cursor, err := TouristsData(td.DB, "tourists").Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("cannot find document in the database %v ", err)
	}

	if err = cursor.All(ctx, &tourists); err != nil {
		return nil, err
	}

	for i := 0; i < len(tourists); i++ {
		tourist := tourists[i]

		bookingSum += len(tourist.BookedTours)
	}

	return &bookingSum, err
}

func (td *AdminDB) SumAllRequestedTours() (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var tourists []model.Tourist

	requestedSum := 0

	cursor, err := TouristsData(td.DB, "tourists").Find(ctx, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("cannot find document in the database %v ", err)
	}

	if err = cursor.All(ctx, &tourists); err != nil {
		return nil, err
	}

	for i := 0; i < len(tourists); i++ {
		tourist := tourists[i]

		requestedSum += len(tourist.RequestTours)
	}

	return &requestedSum, err
}
