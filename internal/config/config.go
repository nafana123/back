package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTSecret          string
	TelegramBotToken   string
	TelegramBotSecret  string
	SteamAPIKey        string
	SteamCallbackURL   string
	SteamRealm         string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	FrontendURL        string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	MailerHost     string
	MailerUser     string
	MailerPassword string
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret:          getEnv("JWT_SECRET", ""),
		TelegramBotToken:   getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramBotSecret:  getEnv("TELEGRAM_BOT_SECRET", ""),
		SteamAPIKey:        getEnv("STEAM_API_KEY", ""),
		SteamCallbackURL:   getEnv("STEAM_CALLBACK_URL", "http://localhost:8080/api/auth/steam/callback"),
		SteamRealm:         getEnv("STEAM_REALM", "http://localhost:8080"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleCallbackURL:  getEnv("GOOGLE_CALLBACK_URL", "http://localhost:8080/api/auth/google/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),

		// Database
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:       getEnv("POSTGRES_DB", "back"),

		// Mailer
		MailerHost:     getEnv("MAILER_HOST", ""),
		MailerUser:     getEnv("MAILER_USER", ""),
		MailerPassword: getEnv("MAILER_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return defaultValue
}
