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

type ServiceRepository struct {
	coll *mongo.Collection
}

type Service struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Duration time.Duration      `bson:"duration"`
}

func NewServiceRepository(client *mongo.Client) *ServiceRepository {
	return &ServiceRepository{
		coll: client.Database("alinea").Collection("services"),
	}
}

func (r *ServiceRepository) FindByID(c context.Context, id string) (core.Service, error) {
	var s Service

	err := r.coll.FindOne(c, bson.D{{Key: "_id", Value: utils.Must(primitive.ObjectIDFromHex(id))}}).Decode(&s)

	if err == mongo.ErrNoDocuments {
		return core.Service{}, nil
	}

	if err != nil {
		return core.Service{}, err
	}

	return *assembleService(s), nil
}

func (r *ServiceRepository) Save(c context.Context, s core.Service) error {
	bToSave := Service{
		ID:       utils.Must(primitive.ObjectIDFromHex(s.ID)),
		Name:     s.Name,
		Duration: s.Duration,
	}

	_, err := r.coll.InsertOne(c, bToSave)

	return err
}

func assembleService(s Service) *core.Service {
	return &core.Service{
		ID:       s.ID.Hex(),
		Name:     s.Name,
		Duration: s.Duration,
	}
}
