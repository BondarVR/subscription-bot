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
