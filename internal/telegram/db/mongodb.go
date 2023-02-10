package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/telegram"
)

type db struct {
	collections *mongo.Collection
	lgr         *logger.LogrusLogger
}

func (d *db) Create(ctx context.Context, user telegram.User) (string, error) {
	result, err := d.collections.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "user is create", err
}

func (d *db) FindOneByID(ctx context.Context, chatId int64) (u telegram.User, err error) {
	filter := bson.M{"chat_id": chatId}

	result := d.collections.FindOne(ctx, filter)
	if result.Err() != nil {
		return u, err
	}

	if err := result.Decode(&u); err != nil {
		return u, err
	}

	return u, nil
}

func (d *db) FindOneByTime(ctx context.Context, time telegram.Time) (u telegram.User, err error) {
	filter := bson.M{"time": time}

	result := d.collections.FindOne(ctx, filter)
	if result.Err() != nil {
		return u, err
	}
	if err := result.Decode(&u); err != nil {
		return u, err
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user telegram.User) error {
	filter := bson.M{"chat_id": user.ChatID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return err
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return err
	}
	delete(updateUserObj, "chat_id")
	delete(updateUserObj, "lon")
	delete(updateUserObj, "lat")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return err
	}

	return nil
}

func NewStorage(database *mongo.Database, collection string, lgr *logger.LogrusLogger) telegram.Storage {
	return &db{
		collections: database.Collection(collection),
		lgr:         lgr,
	}
}
