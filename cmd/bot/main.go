package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/valyala/fasthttp"
	"log"
	"subscription-bot/internal/config"
	"subscription-bot/internal/logger"
	"subscription-bot/internal/repository/mongoDB"
	"subscription-bot/internal/telegram"
	"subscription-bot/internal/telegram/db"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
		lgr.Infof("Server run on http://%s/", addr)
		err := fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {})
		if err != nil {
			lgr.Fatal(err)
		}
	}()

	client, err := mongoDB.NewClient(cfg.MongoUser, cfg.MongoPassword, cfg.MongoHost, cfg.MongoPort, cfg.NameDatabase, context.Background())
	if err != nil {
		lgr.Fatal(err)
	}
	lgr.Infof("DB is start. Name: %s", client.Name())

	storage := db.NewStorage(client, cfg.DbCollections, lgr)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		lgr.Fatal(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, cfg, lgr, storage, wg)
	if err := telegramBot.StartBotAndTicker(); err != nil {
		lgr.Fatal(err)
	}
	wg.Wait()
}
