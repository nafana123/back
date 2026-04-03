package config

import (
	"log"
	"os"
)

type Config struct {
	TelegramBotToken  string
	TelegramBotSecret string
}

func LoadConfig() *Config {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	secret := os.Getenv("TELEGRAM_BOT_SECRET")
	if secret == "" {
		log.Fatal("TELEGRAM_BOT_SECRET не задан")
	}

	return &Config{
		TelegramBotToken:  token,
		TelegramBotSecret: secret,
	}
}
