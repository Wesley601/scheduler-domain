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
		coll: client.Database("alinea").Collection("agendas"),
	}
}

func (r *AgendaRepository) FindByID(id string) (core.Agenda, error) {
	var s Agenda

	err := r.coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: utils.Must(primitive.ObjectIDFromHex(id))}}).Decode(&s)

	if err != nil {
		return core.Agenda{}, err
	}

	return core.Agenda{
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

func (r *AgendaRepository) List() ([]core.Agenda, error) {
	var agendas []core.Agenda

	cur, err := r.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return []core.Agenda{}, err
	}

	for cur.Next(context.Background()) {
		var s Agenda

		err := cur.Decode(&s)
		if err != nil {
			return []core.Agenda{}, err
		}

		agendas = append(agendas, core.Agenda{
			Name: s.Name,
			Slots: func() []core.Slot {
				var slots []core.Slot

				for _, slot := range s.Slots {
					slots = append(slots, utils.Must(core.NewSlot(slot.Weekday, core.SlotTime(slot.StartAt), core.SlotTime(slot.EndsAt))))
				}

				return slots
			}(),
		})
	}

	return agendas, nil
}

func (r *AgendaRepository) Save(s core.Agenda) error {
	var slots []Slot

	for _, slot := range s.Slots {
		slots = append(slots, Slot{
			Weekday: slot.Weekday,
			StartAt: string(slot.StartAt),
			EndsAt:  string(slot.EndsAt),
		})
	}

	sToSave := Agenda{
		ID:    primitive.NewObjectID(),
		Name:  s.Name,
		Slots: slots,
	}

	_, err := r.coll.InsertOne(context.TODO(), sToSave)

	return err
}
