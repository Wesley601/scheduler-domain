package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient() (*mongo.Client, error) {
	return mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
}
