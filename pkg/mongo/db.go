package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017"))
}
