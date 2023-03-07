package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"subscription-bot/internal/config"
	"subscription-bot/internal/models"
	"subscription-bot/internal/telegram"
	"testing"
)

func createMongoDBStorage(t *testing.T) (*mongo.Client, telegram.Storage) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.TestDBURI))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}

	storage := NewStorage(client.Database(cfg.TestDB), cfg.TestDBCollection)

	return client, storage
}

func TestMongoDBStorage_Create(t *testing.T) {
	client, storage := createMongoDBStorage(t)
	defer client.Disconnect(context.Background())
	ctx := context.Background()

	testTable := []struct {
		name string
		in   models.User
	}{
		{
			name: "user is create",
			in: models.User{
				ChatID: 15,
				Lon:    16748,
				Lat:    534,
				Time: models.Time{
					Hour:    "20",
					Minutes: "45",
					Second:  "57",
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			id, err := storage.Create(ctx, tt.in)
			require.NotNil(t, id)
			require.NoError(t, err)
		})
	}
}

func TestMongoDBStorage_FindOneByID(t *testing.T) {
	client, storage := createMongoDBStorage(t)
	defer client.Disconnect(context.Background())
	ctx := context.Background()

	testTable := []struct {
		name      string
		in        int64
		expResult models.User
		expErr    error
	}{
		{
			name: "find user by ID",
			in:   56,
			expResult: models.User{
				ChatID: 56,
				Lon:    10,
				Lat:    50,
				Time: models.Time{
					Hour:    "50",
					Minutes: "15",
					Second:  "15",
				},
			},
			expErr: nil,
		},
		{
			name: "do not find user by ID",
			in:   111111,
			expResult: models.User{
				ChatID: 0,
				Lon:    0,
				Lat:    0,
				Time: models.Time{
					Hour:    "",
					Minutes: "",
					Second:  "",
				},
			},
			expErr: fmt.Errorf("do not find user by ID"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := storage.FindOneByID(ctx, tt.in)

			require.Equal(t, tt.expResult, foundUser)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestMongoDBStorage_FindOneByTime(t *testing.T) {
	client, storage := createMongoDBStorage(t)
	defer client.Disconnect(context.Background())
	ctx := context.Background()

	testTable := []struct {
		name      string
		in        models.Time
		expResult models.User
		expErr    error
	}{
		{
			name: "find user by Time",
			in: models.Time{
				Hour:    "50",
				Minutes: "15",
				Second:  "15",
			},
			expResult: models.User{
				ChatID: 56,
				Lon:    10,
				Lat:    50,
				Time: models.Time{
					Hour:    "50",
					Minutes: "15",
					Second:  "15",
				},
			},
			expErr: nil,
		},
		{
			name: "do not find user by Time",
			in: models.Time{
				Hour:    "00",
				Minutes: "05",
				Second:  "05",
			},
			expResult: models.User{
				ChatID: 0,
				Lon:    0,
				Lat:    0,
				Time: models.Time{
					Hour:    "",
					Minutes: "",
					Second:  "",
				},
			},
			expErr: fmt.Errorf("do not find user by Time"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := storage.FindOneByTime(ctx, tt.in)

			require.Equal(t, tt.expResult, foundUser)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestMongoDBStorage_Update(t *testing.T) {
	client, storage := createMongoDBStorage(t)
	defer client.Disconnect(context.Background())
	ctx := context.Background()

	testTable := []struct {
		name   string
		in     models.User
		expErr error
	}{
		{
			name: "user is update",
			in: models.User{
				ChatID: 14,
				Lon:    1,
				Lat:    5,
				Time: models.Time{
					Hour:    "33",
					Minutes: "33",
					Second:  "3",
				},
			},
			expErr: nil,
		},
		{
			name: "user is not update",
			in: models.User{
				ChatID: 22222,
				Lon:    1,
				Lat:    5,
				Time: models.Time{
					Hour:    "33",
					Minutes: "33",
					Second:  "3",
				},
			},
			expErr: fmt.Errorf("user is not update"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Update(ctx, tt.in)

			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestMongoDBStorage_Delete(t *testing.T) {
	client, storage := createMongoDBStorage(t)
	defer client.Disconnect(context.Background())
	ctx := context.Background()

	testTable := []struct {
		name   string
		in     int64
		count  int64
		expErr error
	}{
		{
			name:   "user is delete",
			in:     15,
			count:  1,
			expErr: nil,
		},
		{
			name:   "user is not delete",
			in:     1111111,
			count:  0,
			expErr: fmt.Errorf("can not delete user by chatId"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			count, err := storage.Delete(ctx, tt.in)
			require.Equal(t, tt.count, count)
			require.Equal(t, tt.expErr, err)
		})
	}
}
