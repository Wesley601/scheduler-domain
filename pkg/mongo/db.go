package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page    int
	PerPage int
}

type Page[T any] struct {
	Data  []T
	Total int64
	Pagination
}

func NewClient(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017"))
}
