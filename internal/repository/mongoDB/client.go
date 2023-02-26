package mongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(username, password, host, port, database string, ctx context.Context) (*mongo.Database, error) {
	var atlasURI string

	if username == "" || password == "" {
		atlasURI = fmt.Sprintf("mongodb://%s:%s/", host, port)
	} else {
		atlasURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/", username, password, host, port)
	}

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
