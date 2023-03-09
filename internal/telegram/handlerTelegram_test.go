package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"subscription-bot/internal/config"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/models"
	mocks "subscription-bot/internal/telegram/mock"
	"testing"
)

var testCtx = context.Background()
var methodForTestHandleTimeFromText = "Update"
var methodForTestHandleStartCommand = "FindOneByID"
var testUser1 = models.User{
	ChatID: 12,
	Lon:    0,
	Lat:    0,
	Time: models.Time{
		Hour:    "20",
		Minutes: "33",
		Second:  "00",
	},
}
var testUser2 = models.User{
	ChatID: 594150834,
	Lon:    12,
	Lat:    15,
	Time: models.Time{
		Hour:    "20",
		Minutes: "33",
		Second:  "00",
	},
}
var testMessage1 = tgbotapi.Message{
	Text: "20:33:00",
	Chat: &tgbotapi.Chat{
		ID: 12,
	},
}
var testMessage2 = tgbotapi.Message{
	Text: "33:33:30",
	Chat: &tgbotapi.Chat{
		ID: 12,
	},
}
var testMessage3 = tgbotapi.Message{
	Chat: &tgbotapi.Chat{
		ID: 594150834,
	},
}

func TestHandleTimeFromText(t *testing.T) {
	testObj := new(mocks.MockStorage)
	cfg, err := config.NewConfig()
	require.NoError(t, err)

	lgr, err := logger.New(cfg)
	require.NoError(t, err)

	bott, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	require.NoError(t, err)

	bot := NewBot(bott, cfg, lgr, testObj)

	testTable := []struct {
		name      string
		inUpdate  models.User
		in        tgbotapi.Message
		expResult string
		method    string
		expErr    error
		resErr    error
	}{
		{
			name:      "User update",
			inUpdate:  testUser1,
			in:        testMessage1,
			expResult: timeText,
			method:    methodForTestHandleTimeFromText,
			expErr:    nil,
			resErr:    nil,
		},
		{
			name:      "Can not parse time",
			inUpdate:  testUser1,
			in:        testMessage2,
			expResult: errTimeText,
			method:    methodForTestHandleTimeFromText,
			expErr:    fmt.Errorf("can not parse time"),
			resErr:    fmt.Errorf("user is not update"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			testObj.On(tt.method, testCtx, tt.inUpdate).Return(tt.resErr)
			result, err := bot.handleTimeFromText(&tt.in)
			require.Equal(t, tt.expErr, err)
			require.Equal(t, tt.expResult, result)
		})
	}
}

func TestHandleStartCommand(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	cfg, err := config.NewConfig()
	require.NoError(t, err)

	lgr, err := logger.New(cfg)
	require.NoError(t, err)

	bott, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	require.NoError(t, err)

	bot := NewBot(bott, cfg, lgr, mockStorage)

	testTable := []struct {
		name             string
		inFind           int64
		expResultFind    models.User
		in               tgbotapi.Message
		methodForStorage string
		expErrForStorage error
	}{
		{name: "User is create",
			inFind:           594150834,
			expResultFind:    testUser2,
			in:               testMessage3,
			methodForStorage: methodForTestHandleStartCommand,
			expErrForStorage: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On(tt.methodForStorage, testCtx, tt.inFind).Return(tt.expResultFind, tt.expErrForStorage)
			err := bot.handleStartCommand(&tt.in)
			require.NoError(t, err)
		})
	}
}
