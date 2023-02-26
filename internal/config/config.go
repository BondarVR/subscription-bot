package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	ApiWeather    string `env:"API_WEATHER"`

	LogServer   string `env:"LOG_SERVER"`
	Loglevel    string `env:"LOG_LEVEL"`
	ServiceName string `env:"SERVICE_NAME"`

	ServerHost string `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort string `env:"SERVER_PORT" envDefault:"8000"`

	MongoHost     string `env:"MONGO_HOST" envDefault:"localhost"`
	MongoPort     string `env:"MONGO_PORT" envDefault:"8010"`
	MongoUser     string `env:"MONGO_USER"`
	NameDatabase  string `env:"MONGO_DB"`
	MongoPassword string `env:"MONGO_PASSWORD"`

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
