package query

import (
	"context"
	"fmt"
	"time"
	"travas_admin/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (td *AdminDB) FindAllTourists() (tourists []model.Tourist, err error) {
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
