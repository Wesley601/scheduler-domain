package mongo

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScheduleRepository struct {
	coll *mongo.Collection
}

type Slot struct {
	Weekday time.Weekday
	StartAt string
	EndsAt  string
}

type Schedule struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Slots []Slot             `bson:"slots"`
}

func NewScheduleRepository(client *mongo.Client) *ScheduleRepository {
	return &ScheduleRepository{
		coll: client.Database("alinea").Collection("schedules"),
	}
}

func (r *ScheduleRepository) FindByID(id string) (core.Schedule, error) {
	var s Schedule

	err := r.coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: utils.Must(primitive.ObjectIDFromHex(id))}}).Decode(&s)

	if err != nil {
		return core.Schedule{}, err
	}

	return core.Schedule{
		Name: s.Name,
		Slots: func() []core.Slot {
			var slots []core.Slot

			for _, slot := range s.Slots {
				slots = append(slots, utils.Must(core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))))
			}

			return slots
		}(),
	}, nil
}

func (r *ScheduleRepository) Save(s core.Schedule) error {
	var slots []Slot

	for _, slot := range s.Slots {
		slots = append(slots, Slot{
			Weekday: slot.Weekday,
			StartAt: string(slot.StartAt),
			EndsAt:  string(slot.EndsAt),
		})
	}

	sToSave := Schedule{
		ID:    primitive.NewObjectID(),
		Name:  s.Name,
		Slots: slots,
	}

	_, err := r.coll.InsertOne(context.TODO(), sToSave)

	return err
}
