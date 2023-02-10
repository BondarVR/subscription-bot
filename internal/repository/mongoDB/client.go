package mongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(username, password, database string, ctx context.Context) (*mongo.Database, error) {
	atlasURI := fmt.Sprintf("mongodb+srv://%s:%s@descriptionbot.m0rvk0o.mongodb.net/?retryWrites=true&w=majority", username, password)
	clientOptions := options.Client().ApplyURI(atlasURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(database), nil
}
