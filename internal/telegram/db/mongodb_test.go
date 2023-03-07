package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"subscription-bot/internal/config"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/models"
	"subscription-bot/internal/telegram"
	"testing"
)

var testCtx = context.Background()
var testTime1 = models.Time{
	Hour:    "50",
	Minutes: "15",
	Second:  "15",
}
var testTime2 = models.Time{
	Hour:    "00",
	Minutes: "05",
	Second:  "05",
}
var testUser1 = models.User{
	ChatID: 15,
	Lon:    16748,
	Lat:    534,
	Time: models.Time{
		Hour:    "20",
		Minutes: "45",
		Second:  "57",
	},
}
var testUser2 = models.User{
	ChatID: 56,
	Lon:    10,
	Lat:    50,
	Time: models.Time{
		Hour:    "50",
		Minutes: "15",
		Second:  "15",
	},
}
var testUser3 = models.User{
	ChatID: 0,
	Lon:    0,
	Lat:    0,
	Time: models.Time{
		Hour:    "",
		Minutes: "",
		Second:  "",
	},
}
var testUser4 = models.User{
	ChatID: 14,
	Lon:    1,
	Lat:    5,
	Time: models.Time{
		Hour:    "33",
		Minutes: "33",
		Second:  "3",
	},
}
var testUser5 = models.User{
	ChatID: 22222,
	Lon:    1,
	Lat:    5,
	Time: models.Time{
		Hour:    "33",
		Minutes: "33",
		Second:  "3",
	},
}

func createTestMongoDBStorage(t *testing.T) (*mongo.Client, telegram.Storage) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)

	lgr, err := logger.New(cfg)
	require.NoError(t, err)

	client, err := mongo.Connect(testCtx, options.Client().ApplyURI(cfg.TestDBURI))
	require.NoError(t, err)

	storage := NewStorage(client.Database(cfg.TestDB), cfg.TestDBCollection, lgr)

	return client, storage
}

func TestMongoDBStorage_Create(t *testing.T) {
	client, storage := createTestMongoDBStorage(t)
	t.Cleanup(func() {
		err := client.Disconnect(testCtx)
		require.NoError(t, err)
	})

	testTable := []struct {
		name   string
		in     models.User
		chatId int64
	}{
		{
			name:   "user is create",
			in:     testUser1,
			chatId: 15,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			id, err := storage.Create(testCtx, tt.in)
			require.NotNil(t, id)
			require.NoError(t, err)

			user, err := storage.FindOneByID(testCtx, tt.chatId)
			require.NoError(t, err)
			require.Equal(t, tt.in, user)
		})
	}
}

func TestMongoDBStorage_FindOneByID(t *testing.T) {
	client, storage := createTestMongoDBStorage(t)
	t.Cleanup(func() {
		err := client.Disconnect(testCtx)
		require.NoError(t, err)
	})

	testTable := []struct {
		name      string
		in        int64
		expResult models.User
		expErr    error
	}{
		{
			name:      "find user by ID",
			in:        56,
			expResult: testUser2,
			expErr:    nil,
		},
		{
			name:      "do not find user by ID",
			in:        111111,
			expResult: testUser3,
			expErr:    fmt.Errorf("do not find user by ID"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := storage.FindOneByID(testCtx, tt.in)
			require.Equal(t, tt.expResult, foundUser)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestMongoDBStorage_FindOneByTime(t *testing.T) {
	client, storage := createTestMongoDBStorage(t)
	t.Cleanup(func() {
		err := client.Disconnect(testCtx)
		require.NoError(t, err)
	})

	testTable := []struct {
		name      string
		in        models.Time
		expResult models.User
		expErr    error
	}{
		{
			name:      "find user by Time",
			in:        testTime1,
			expResult: testUser2,
			expErr:    nil,
		},
		{
			name:      "do not find user by Time",
			in:        testTime2,
			expResult: testUser3,
			expErr:    fmt.Errorf("do not find user by Time"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := storage.FindOneByTime(testCtx, tt.in)
			require.Equal(t, tt.expResult, foundUser)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestMongoDBStorage_Update(t *testing.T) {
	client, storage := createTestMongoDBStorage(t)
	t.Cleanup(func() {
		err := client.Disconnect(testCtx)
		require.NoError(t, err)
	})

	testTable := []struct {
		name      string
		in        models.User
		expErr    error
		chatId    int64
		resErr    error
		expResult models.User
	}{
		{
			name:      "user is update",
			in:        testUser4,
			expErr:    nil,
			chatId:    14,
			resErr:    nil,
			expResult: testUser4,
		},
		{
			name:      "user is not update",
			in:        testUser5,
			expErr:    fmt.Errorf("user is not update"),
			chatId:    22222,
			resErr:    fmt.Errorf("do not find user by ID"),
			expResult: testUser3,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Update(testCtx, tt.in)
			require.Equal(t, tt.expErr, err)

			user, err := storage.FindOneByID(testCtx, tt.chatId)
			require.Equal(t, tt.expResult, user)
			require.Equal(t, tt.resErr, err)
		})
	}
}

func TestMongoDBStorage_Delete(t *testing.T) {
	client, storage := createTestMongoDBStorage(t)
	t.Cleanup(func() {
		err := client.Disconnect(testCtx)
		require.NoError(t, err)
	})

	testTable := []struct {
		name      string
		chatId    int64
		count     int64
		expErr    error
		resErr    error
		expResult models.User
	}{
		{
			name:      "user is delete",
			chatId:    15,
			count:     1,
			expErr:    nil,
			resErr:    fmt.Errorf("do not find user by ID"),
			expResult: testUser3,
		},
		{
			name:      "user is not delete",
			chatId:    1111111,
			count:     0,
			expErr:    fmt.Errorf("can not delete user by chatId"),
			resErr:    fmt.Errorf("do not find user by ID"),
			expResult: testUser3,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			count, err := storage.Delete(testCtx, tt.chatId)
			require.Equal(t, tt.count, count)
			require.Equal(t, tt.expErr, err)

			user, err := storage.FindOneByID(testCtx, tt.chatId)
			require.Equal(t, tt.expResult, user)
			require.Equal(t, tt.resErr, err)
		})
	}
}
