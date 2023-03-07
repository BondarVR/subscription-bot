package models

type User struct {
	ChatID int64   `json:"chat_id" bson:"chat_id"`
	Lon    float64 `json:"lon" bson:"lon"`
	Lat    float64 `json:"lat" bson:"lat"`
	Time   Time    `json:"time" bson:"time"`
}

type Time struct {
	Hour    string `json:"hour" bson:"hour"`
	Minutes string `json:"minutes" bson:"minutes"`
	Second  string `json:"second" bson:"second"`
}

type Message struct {
	Coord      Coord   `json:"coord"`
	Weather    Weather `json:"weather"`
	Base       string  `json:"base"`
	Main       Main    `json:"main"`
	Visibility int     `json:"visibility"`
	Wind       Wind    `json:"wind"`
	Clouds     Clouds  `json:"clouds"`
	Dt         int     `json:"dt"`
	Sys        Sys     `json:"sys"`
	Timezone   int     `json:"timezone"`
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Cod        int     `json:"cod"`
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
