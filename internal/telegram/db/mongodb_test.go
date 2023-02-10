package db

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"subscription-bot/internal/telegram"
	mock_telegram "subscription-bot/internal/telegram/mocks"
	"testing"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	in := telegram.User{
		ChatID: 1,
		Lon:    35,
		Lat:    46,
		Time: telegram.Time{
			Hour:    "21",
			Minutes: "43",
			Second:  "56",
		},
	}
	mockResult := "user is create"

	repo := mock_telegram.NewMockStorage(ctrl)

	repo.EXPECT().Create(ctx, in).Return(mockResult, nil).Times(1)
	create, err := repo.Create(ctx, in)
	require.NoError(t, err)
	require.Equal(t, mockResult, create)
}

func TestFindOneByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	in := int64(25)
	mockResult := telegram.User{
		ChatID: 25,
		Lon:    2345,
		Lat:    4638972,
		Time: telegram.Time{
			Hour:    "22",
			Minutes: "45",
			Second:  "23",
		},
	}

	repo := mock_telegram.NewMockStorage(ctrl)
	repo.EXPECT().FindOneByID(ctx, in).Return(mockResult, nil).Times(1)
	user, err := repo.FindOneByID(ctx, in)
	require.NoError(t, err)
	require.Equal(t, mockResult, user)
}

func TestFindOneByTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	in := telegram.Time{
		Hour:    "19",
		Minutes: "11",
		Second:  "01",
	}
	mockResult := telegram.User{
		ChatID: 1,
		Lon:    29,
		Lat:    45,
		Time: telegram.Time{
			Hour:    "19",
			Minutes: "11",
			Second:  "01",
		},
	}

	repo := mock_telegram.NewMockStorage(ctrl)
	repo.EXPECT().FindOneByTime(ctx, in).Return(mockResult, nil).Times(1)
	user, err := repo.FindOneByTime(ctx, in)
	require.NoError(t, err)
	require.Equal(t, mockResult, user)
}
