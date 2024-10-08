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

type BlockRepository struct {
	coll *mongo.Collection
}

type Block struct {
	ID      primitive.ObjectID `bson:"_id"`
	Weekday time.Weekday       `bson:"weekday"`
	From    primitive.DateTime `bson:"from"`
	To      primitive.DateTime `bson:"to"`
}

func NewBlockRepository(client *mongo.Client) *BlockRepository {
	return &BlockRepository{
		coll: client.Database("base").Collection("blocks"),
	}
}

func (r *BlockRepository) IsAvailable(c context.Context, w core.Window) (bool, error) {
	var b Block

	err := r.coll.FindOne(c, bson.D{
		{Key: "from", Value: w.From},
		{Key: "to", Value: w.To},
	}).Decode(&b)

	if err == mongo.ErrNoDocuments {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (r *BlockRepository) Save(c context.Context, s core.Block) error {
	bToSave := Block{
		ID:      utils.Must(primitive.ObjectIDFromHex(s.ID)),
		Weekday: s.Weekday,
		From:    primitive.NewDateTimeFromTime(s.Window.From),
		To:      primitive.NewDateTimeFromTime(s.Window.To),
	}

	_, err := r.coll.InsertOne(c, bToSave)

	return err
}
