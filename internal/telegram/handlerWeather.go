package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"time"
)

func (b *Bot) handleParamForWeather(message *Message) string {
	description := message.Weather[0].Description
	mainTemp := b.handleTemperature(message.Main.Temp)
	feelsLike := b.handleTemperature(message.Main.FeelsLike)
	pressure := b.handlePressure(message.Main.Pressure)
	humidity := message.Main.Humidity
	windSpeed := message.Wind.Speed
	clouds := message.Clouds.All
	sunrise := b.handleDate(message.Sys.Sunrise)
	sunset := b.handleDate(message.Sys.Sunset)
	nameLocation := message.Name

	text := fmt.Sprintf("Локація: %s\n"+
		"Стан погоди: %s\n"+
		"Температура: %.1f градусів\n"+
		"Відчувається як: %.1f градусів\n"+
		"Атмосферний тиск: %.1f мм.рт.ст.\n"+
		"Вологість повітря у відсотках: %d\n"+
		"Швидкість вітру : %.2f м/с\n"+
		"Хмарність у відсотках: %d\n"+
		"Світанок: %s\n"+
		"Закат: %s\n",
		nameLocation,
		description,
		mainTemp,
		feelsLike,
		pressure,
		humidity,
		windSpeed,
		clouds,
		sunrise,
		sunset)

	return text
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleTemperature(temp float32) float32 {
	handleTemp := temp - 273.15
	return handleTemp
}

func (b *Bot) handleDate(date int) string {
	dateConvert := int64(date)
	timeT := time.Unix(dateConvert, 0)
	timeF := timeT.Format("15:04:05")
	return timeF
}

func (b *Bot) handlePressure(pressure int) float32 {
	formatPressure := float32(pressure) / 1.333
	return formatPressure
}

func (b *Bot) unmarshalJSON(response *http.Response) (*Message, error) {
	var v Message
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
