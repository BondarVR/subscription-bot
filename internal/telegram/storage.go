package telegram

import (
	"context"
	"subscription-bot/internal/models"
)

type Storage interface {
	Create(ctx context.Context, user models.User) (string, error)
	FindOneByID(ctx context.Context, chatId int64) (u models.User, err error)
	FindOneByTime(ctx context.Context, time models.Time) (u models.User, err error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, chatId int64) (int64, error)
}
