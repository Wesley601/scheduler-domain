package mongo

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository struct {
	coll *mongo.Collection
}

type Booking struct {
	ID      primitive.ObjectID `bson:"_id"`
	StartAt time.Time          `bson:"startAt"`
	EndsAt  time.Time          `bson:"endsAt"`
}

func NewBookingRepository(client *mongo.Client) *BookingRepository {
	return &BookingRepository{
		coll: client.Database("alinea").Collection("bookings"),
	}
}

func (r *BookingRepository) IsAvailable(w core.Window) (bool, error) {
	var b Booking

	err := r.coll.FindOne(context.TODO(), bson.D{
		{Key: "startAt", Value: w.From},
		{Key: "endsAt", Value: w.To},
	}).Decode(&b)

	if err == mongo.ErrNoDocuments {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (r *BookingRepository) Save(b core.Booking) error {
	bToSave := Booking{
		ID:      primitive.NewObjectID(),
		StartAt: b.Window.From,
		EndsAt:  b.Window.To,
	}

	_, err := r.coll.InsertOne(context.TODO(), bToSave)

	return err
}
