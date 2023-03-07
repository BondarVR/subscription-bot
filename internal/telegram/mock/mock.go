package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"subscription-bot/internal/models"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Create(ctx context.Context, user models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockStorage) FindOneByID(ctx context.Context, chatId int64) (u models.User, err error) {
	args := m.Called(ctx, chatId)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) FindOneByTime(ctx context.Context, time models.Time) (u models.User, err error) {
	args := m.Called(ctx, time)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockStorage) Update(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStorage) Delete(ctx context.Context, chatId int64) (int64, error) {
	args := m.Called(ctx, chatId)
	return args.Get(0).(int64), args.Error(1)
}
