package config

import (
	"sync"

	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramBotToken string `env:"BOTTOKEN" env-required:true`
	DatabaseURL      string `env:"DATABASEURL" env-default:"postgres://postgres:postgres@localhost:5432/chsuBot?sslmode=disable"`
	AdminId          int    `env:"ADMIN" env-required:true`
	TypeStart        string `env:"TYPESTART" env-default:"long-polling"`
	WebhookURL       string `env:"WEBHOOKURL"`
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.New()
		logger.Info("Read app configuration")
		cfg = &Config{}
		if err := cleanenv.ReadEnv(cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return cfg
}
