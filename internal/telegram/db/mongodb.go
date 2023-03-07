package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/models"
	"subscription-bot/internal/telegram"
)

type db struct {
	collections *mongo.Collection
	lgr         *logger.LogrusLogger
}

func NewStorage(database *mongo.Database, collection string) telegram.Storage {
	return &db{
		collections: database.Collection(collection),
	}
}

func (d *db) Create(ctx context.Context, user models.User) (string, error) {
	result, err := d.collections.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to conver objectid to hex")
}

func (d *db) FindOneByID(ctx context.Context, chatId int64) (u models.User, err error) {
	filter := bson.M{"chat_id": chatId}

	result := d.collections.FindOne(ctx, filter)
	if result.Err() != nil {
		return u, fmt.Errorf("do not find user by ID")
	}

	if err := result.Decode(&u); err != nil {
		return u, err
	}

	return u, nil
}

func (d *db) FindOneByTime(ctx context.Context, time models.Time) (u models.User, err error) {
	filter := bson.M{"time": time}

	result := d.collections.FindOne(ctx, filter)
	if result.Err() != nil {
		return u, fmt.Errorf("do not find user by Time")
	}
	if err := result.Decode(&u); err != nil {
		return u, err
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user models.User) error {
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
		return fmt.Errorf("user is not update")
	}

	return nil
}

func (d *db) Delete(ctx context.Context, chatId int64) (int64, error) {
	filter := bson.M{"chat_id": chatId}
	result, err := d.collections.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("can not delete user by chatId")
	}
	return result.DeletedCount, nil
}
