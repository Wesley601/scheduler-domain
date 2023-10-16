package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"alinea.com/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
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
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	wg.Add(1)
	go func() {
		utils.LoadingMessage(ctx, "starting to connect on mongodb...")
		wg.Done()
	}()

	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}

	wg.Wait()

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}
