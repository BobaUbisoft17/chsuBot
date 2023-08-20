package config

import (
	"sync"

	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramBotToken string `env:"BOTTOKEN" env-required:true`
	DatabaseURL      string `env:"DATABASEURL" env-default:"postgres://postgres:postgres@localhost:5432/chsuBot?sslmode=disable"`
	Admin            int    `env:"ADMIN"`
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.New()
		logger.Info("Чтение конфигурации приложения")
		cfg = &Config{}
		if err := cleanenv.ReadEnv(cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return cfg
}