package config

import (
	"os"
)

type Config struct {
	JWTSecret         string
	TelegramBotToken  string
	TelegramBotSecret string
	SteamAPIKey       string
	SteamCallbackURL  string
	SteamRealm        string
	FrontendURL       string
	
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret:         getEnv("JWT_SECRET", "default-secret-key"),
		TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramBotSecret: getEnv("TELEGRAM_BOT_SECRET", ""),
		SteamAPIKey:       getEnv("STEAM_API_KEY", ""),
		SteamCallbackURL:  getEnv("STEAM_CALLBACK_URL", "http://localhost:8080/api/auth/steam/callback"),
		SteamRealm:        getEnv("STEAM_REALM", "http://localhost:8080"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:5173"),
		
		// Database defaults
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:       getEnv("POSTGRES_DB", "back"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}