package telegram

import (
	"fmt"
	"net/http"
)

func (b *Bot) GetWeatherInfo(lon, lat float64, chatID int64) error {
	if lon == 0 || lat == 0 || b.cfg.ApiWeather == "" {
		return fmt.Errorf("link creation data is empty")
	}
	requestURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&lang=ua&appid=%s", lat, lon, b.cfg.ApiWeather)

	resp, err := http.Get(requestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	value, err := b.unmarshalJSON(resp)
	if err != nil {
		return err
	}

	text, err := b.handleParamForWeather(value)
	if err != nil {
		return err
	}

	if err := b.sendMessage(chatID, text); err != nil {
		return err
	}

	return nil
}
