package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"subscription-bot/internal/config"
	"subscription-bot/internal/logger"
	"sync"
	"time"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	cfg     *config.Config
	lgr     *logger.LogrusLogger
	storage Storage
	wg      sync.WaitGroup
}

func NewBot(bot *tgbotapi.BotAPI, cfg *config.Config, lgr *logger.LogrusLogger, storage Storage, wg sync.WaitGroup) *Bot {
	return &Bot{bot: bot, cfg: cfg, lgr: lgr, storage: storage, wg: wg}
}

func (b *Bot) StartBotAndTicker() error {
	b.lgr.Infof("Authorized on account %s", b.bot.Self.UserName)
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		ticker := time.NewTicker(45 * time.Second)
		for _ = range ticker.C {
			b.lgr.Infof("tick...")
			realTime := TimeForTicker()
			user, err := b.storage.FindOneByTime(context.Background(), realTime)
			if err != nil {
				return
			}
			if user.ChatID != 0 {
				err = b.GetWeatherInfo(user.Lon, user.Lat, user.ChatID)
				if err != nil {
					return
				}
				b.lgr.Infof("200:OK (send message)")
			}
		}
	}()
	updates := b.botUpdateChannel()
	if err := b.handleUpdate(updates); err != nil {
		return err
	}
	b.wg.Wait()
	return nil
}

func (b *Bot) botUpdateChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdate(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				return err
			}
			continue
		}
		if update.Message.Location != nil {
			if err := b.handleLocation(update.Message); err != nil {
				return err
			}
			continue
		}
		if err := b.handleText(update.Message); err != nil {
			return nil
		}
	}
	return nil
}
