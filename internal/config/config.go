package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	LogServer     string `env:"LOG_SERVER"`
	Loglevel      string `env:"LOG_LEVEL"`
	ServiceName   string `env:"SERVICE_NAME"`
	ApiWeather    string `env:"API_WEATHER"`
	MongoUser     string `env:"MONGO_USER"`
	MongoPassword string `env:"MONGO_PASSWORD"`
	DbName        string `env:"DB_NAME"`
	DbCollections string `env:"DB_COLLECTIONS"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
