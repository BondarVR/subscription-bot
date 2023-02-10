package telegram

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindOneByID(ctx context.Context, chatId int64) (u User, err error)
	FindOneByTime(ctx context.Context, time Time) (u User, err error)
	Update(ctx context.Context, user User) error
}
