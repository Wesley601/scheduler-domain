package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"wesley601.com/internal/core"
	"wesley601.com/pkg/utils"
)

type AgendaRepository struct {
	coll *mongo.Collection
}

type Slot struct {
	Weekday time.Weekday
	StartAt string
	EndsAt  string
}

type Agenda struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Slots []Slot             `bson:"slots"`
}

func NewAgendaRepository(client *mongo.Client) *AgendaRepository {
	return &AgendaRepository{
		coll: client.Database("base").Collection("agendas"),
	}
}

func (r *AgendaRepository) FindByID(c context.Context, id string) (core.Agenda, error) {
	var s Agenda

	err := r.coll.FindOne(c, bson.D{{Key: "_id", Value: utils.Must(primitive.ObjectIDFromHex(id))}}).Decode(&s)

	if err != nil {
		return core.Agenda{}, err
	}

	return assembleAgenda(s), nil
}

func (r *AgendaRepository) List(c context.Context) ([]core.Agenda, error) {
	var agendas []core.Agenda

	cur, err := r.coll.Find(c, bson.D{})
	if err != nil {
		return []core.Agenda{}, err
	}

	for cur.Next(c) {
		var s Agenda

		err := cur.Decode(&s)
		if err != nil {
			return []core.Agenda{}, err
		}

		agendas = append(agendas, assembleAgenda(s))
	}

	return agendas, nil
}

func (r *AgendaRepository) Save(c context.Context, s core.Agenda) error {
	var slots []Slot

	for _, slot := range s.Slots {
		slots = append(slots, Slot{
			Weekday: slot.Weekday,
			StartAt: string(slot.StartAt),
			EndsAt:  string(slot.EndsAt),
		})
	}

	sToSave := Agenda{
		ID:    utils.Must(primitive.ObjectIDFromHex(s.ID)),
		Name:  s.Name,
		Slots: slots,
	}

	_, err := r.coll.InsertOne(c, sToSave)

	return err
}

func assembleAgenda(s Agenda) core.Agenda {
	return core.Agenda{
		ID:   s.ID.Hex(),
		Name: s.Name,
		Slots: func() []core.Slot {
			var slots []core.Slot

			for _, slot := range s.Slots {
				slots = append(slots, utils.Must(core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))))
			}

			return slots
		}(),
	}

}
