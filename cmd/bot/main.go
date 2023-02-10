package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"subscription-bot/internal/config"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/repository/mongoDB"
	"subscription-bot/internal/telegram"
	"subscription-bot/internal/telegram/db"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	lgr, err := logger.New(logger.Config{
		LogServer:   cfg.LogServer,
		LogLevel:    cfg.Loglevel,
		ServiceName: cfg.ServiceName,
	})
	if err != nil {
		lgr.Fatal(err)
	}

	client, err := mongoDB.NewClient(cfg.MongoUser, cfg.MongoPassword, cfg.DbName, context.Background())
	if err != nil {
		lgr.Fatal(err)
	}

	storage := db.NewStorage(client, cfg.DbCollections, lgr)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		lgr.Fatal(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, cfg, lgr, storage)
	if err := telegramBot.StartBotAndTicker(); err != nil {
		lgr.Fatal(err)
	}
}
