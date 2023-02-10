package telegram

import (
	"fmt"
	"net/http"
)

type Message struct {
	Coord      `json:"coord"`
	Weather    `json:"weather"`
	Base       string `json:"base"`
	Main       `json:"main"`
	Visibility int `json:"visibility"`
	Wind       `json:"wind"`
	Clouds     `json:"clouds"`
	Dt         int `json:"dt"`
	Sys        `json:"sys"`
	Timezone   int    `json:"timezone"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Cod        int    `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather []struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float32 `json:"temp"`
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type Clouds struct {
	All int `json:"all"`
}

type Wind struct {
	Speed float32 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float32 `json:"gust"`
}

type Sys struct {
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

func (b *Bot) GetWeatherInfo(lon, lat float64, chatID int64) error {
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

	text := b.handleParamForWeather(value)

	if err := b.sendMessage(chatID, text); err != nil {
		return err
	}

	return nil
}
