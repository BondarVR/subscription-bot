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

func TestHandleTimeFromText(t *testing.T) {
	testObj := new(mocks.MockStorage)
	ctx := context.Background()
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
	}{
		{name: "User update",
			inUpdate: models.User{
				ChatID: 12,
				Lon:    0,
				Lat:    0,
				Time: models.Time{
					Hour:    "20",
					Minutes: "33",
					Second:  "00",
				},
			},
			in: tgbotapi.Message{
				Text: "20:33:30",
				Chat: &tgbotapi.Chat{
					ID: 12,
				},
			},
			expResult: timeText,
			method:    "Update",
			expErr:    nil,
		},
		{name: "Can not parse time",
			inUpdate: models.User{
				ChatID: 12,
				Lon:    0,
				Lat:    0,
				Time: models.Time{
					Hour:    "20",
					Minutes: "33",
					Second:  "00",
				},
			},
			in: tgbotapi.Message{
				Text: "33:33:30",
				Chat: &tgbotapi.Chat{
					ID: 12,
				},
			},
			expResult: errTimeText,
			method:    "Update",
			expErr:    fmt.Errorf("can not parse time"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			testObj.On(tt.method, ctx, tt.inUpdate).Return(tt.expErr)
			result, err := bot.handleTimeFromText(&tt.in)
			testObj.AssertExpectations(t)

			require.Equal(t, err, tt.expErr)
			require.Equal(t, tt.expResult, result)
		})
	}
}

func TestHandleStartCommand(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	ctx := context.Background()
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
			inFind: 594150834,
			expResultFind: models.User{
				ChatID: 594150834,
				Lon:    12,
				Lat:    15,
				Time: models.Time{
					Hour:    "20",
					Minutes: "33",
					Second:  "00",
				},
			},
			in: tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 594150834,
				},
			},
			methodForStorage: "FindOneByID",
			expErrForStorage: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage.On(tt.methodForStorage, ctx, tt.inFind).Return(tt.expResultFind, tt.expErrForStorage)
			err := bot.handleStartCommand(&tt.in)
			mockStorage.AssertExpectations(t)

			require.NoError(t, err)
		})
	}
}
