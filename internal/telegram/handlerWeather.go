package telegram

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
	"subscription-bot/internal/models"
	"time"
)

func (b *Bot) handleParamForWeather(message *models.Message) (string, error) {
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

	err := validateWeatherData(nameLocation, description, mainTemp, feelsLike, pressure, humidity, windSpeed, clouds, sunrise, sunset)
	if err != nil {
		return "", err
	}

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

	return text, nil
}

func validateWeatherData(nameLocation, description string, mainTemp, feelsLike, pressure float32, humidity int, windSpeed float32, clouds int, sunrise, sunset string) error {
	var errStrings []string
	if nameLocation == "" {
		errStrings = append(errStrings, "Name location is empty")
	}
	if description == "" {
		errStrings = append(errStrings, "Weather description is empty")
	}
	if mainTemp == 0 {
		errStrings = append(errStrings, "Main temperature is empty")
	}
	if feelsLike == 0 {
		errStrings = append(errStrings, "Feels like temperature is empty")
	}
	if pressure == 0 {
		errStrings = append(errStrings, "Atmospheric pressure is empty")
	}
	if humidity == 0 {
		errStrings = append(errStrings, "Humidity is empty")
	}
	if windSpeed == 0 {
		errStrings = append(errStrings, "Wind speed is empty")
	}
	if clouds == 0 {
		errStrings = append(errStrings, "Clouds is empty")
	}
	if sunrise == "" {
		errStrings = append(errStrings, "Sunrise time is empty")
	}
	if sunset == "" {
		errStrings = append(errStrings, "Sunset time is empty")
	}

	if len(errStrings) > 0 {
		return fmt.Errorf("Invalid weather data: %s", strings.Join(errStrings, ", "))
	}
	return nil
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

func (b *Bot) unmarshalJSON(response *http.Response) (*models.Message, error) {
	var v models.Message
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
