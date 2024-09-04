package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wesley601.com/internal/core"
	"wesley601.com/pkg/utils"
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

type ListFilter struct {
	Q string
	Pagination
}

type ServicePage Page[core.Service]

func (r *ServiceRepository) List(c context.Context, f ListFilter) (ServicePage, error) {
	var services []core.Service
	page := ServicePage{
		Pagination: Pagination{
			Page:    f.Page,
			PerPage: f.PerPage,
		},
	}

	total, err := r.coll.CountDocuments(c, bson.D{})
	if err != nil {
		return page, err
	}
	page.Total = total

	skip := int64(f.Page*f.PerPage - f.PerPage)
	l := int64(f.PerPage)

	cursor, err := r.coll.Find(c, bson.D{{Key: "name", Value: primitive.Regex{
		Pattern: f.Q,
		Options: "i",
	}}}, &options.FindOptions{Limit: &l, Skip: &skip})
	if err != nil {
		return page, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var service Service

		err := cursor.Decode(&service)
		if err != nil {
			return page, err
		}

		services = append(services, core.Service{
			ID:       service.ID.Hex(),
			Name:     service.Name,
			Duration: service.Duration,
		})
	}

	page.Data = services

	return page, nil
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
