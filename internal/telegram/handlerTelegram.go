package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"subscription-bot/internal/models"
	"time"
)

const (
	commandStart = "start"
	errCommand   = "Invalid type of command."
	firstButton  = "Реєстрація"
	startText    = "Цей бот може надавати Вам інформацію про погодні умови:\n\n" +
		"Натискайте 'Реєстрація' щоб почати користування ботом!"
	registerUser    = "Ви вже налаштували бота для відправки Вам повідомлень!\n"
	instructionText = "Інструкція як отримати підписку на сповіщення про погоду. \n\n" +
		"Спочатку потрібно визначитись із Вашою локацією, для цього:\n" +
		"1. Натисніть на скріпку праворуч внизу екрану\n" +
		"2. Оберіть у нижньому меню ʼГеопозицияʼ\n" +
		"3. Оберіть перший пункт ʼОтправить свою геопозициюʼ\n"
	geoText = "Ваша локація збережена, тепер потрібно визначитись із часом відправки повідомлень.\n\n" +
		"Ведіть часу форматі: 'години:хвилини:секунди'.\n\n" +
		"Наприклад: 08:30:00, або 20:30:00"
	errTimeText = "Ви вказали час у неправильному форматі, спробуйте ще раз.\n" +
		"Ведіть часу форматі: 'години:хвилини:секунди'.\n\n" +
		"Наприклад: 08:30:00, або 20:30:00"
	timeText = "Обраний Вами час збережено.\n" +
		"Налаштування закінчено. \n\n" +
		"Тепер Вам щоденно будуть приходити сповіщення про погоду у вибраний час.\n\n" +
		"Хай щастить!"
	layoutOfTime = "15:04:05"
)

func (b *Bot) handleText(message *tgbotapi.Message) error {
	switch message.Text {
	case firstButton:
		paramKeyboard := tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		msg := tgbotapi.NewMessage(message.Chat.ID, instructionText)
		msg.ReplyMarkup = paramKeyboard
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	default:
		text, err := b.handleTimeFromText(message)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		if err := b.handleStartCommand(message); err != nil {
			return err
		}
		return nil
	default:
		return errors.New(errCommand)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	checkUser, err := b.storage.FindOneByID(context.Background(), message.Chat.ID)
	if err != nil {
		return err
	}

	if checkUser.ChatID == message.Chat.ID {
		msg := tgbotapi.NewMessage(message.Chat.ID, registerUser)
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
		return nil
	}

	paramKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(firstButton),
		),
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, startText)
	msg.ReplyMarkup = paramKeyboard
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleLocation(message *tgbotapi.Message) error {
	user := models.User{
		ChatID: int64(message.Chat.ID),
		Lon:    message.Location.Longitude,
		Lat:    message.Location.Latitude,
		Time: models.Time{
			Hour:    "00",
			Minutes: "00",
			Second:  "00",
		},
	}
	result, err := b.storage.Create(context.Background(), user)
	if err != nil {
		return err
	}

	b.lgr.Info(result)

	msg := tgbotapi.NewMessage(message.Chat.ID, geoText)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleTimeFromText(message *tgbotapi.Message) (string, error) {
	_, err := time.Parse(layoutOfTime, message.Text)
	if err != nil {
		return errTimeText, fmt.Errorf("can not parse time")
	} else {
		time := ParseTime(message.Text)
		user := models.User{
			ChatID: int64(message.Chat.ID),
			Time:   time,
		}
		err = b.storage.Update(context.Background(), user)
		if err != nil {
			return "", err
		}
	}
	return timeText, nil
}
