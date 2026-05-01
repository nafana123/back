package main

import (
	"back/internal/app"
	"log"
	"net/http"

	"back/internal/config"
	"back/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env файл не загружен")
	}

	logg := logger.New()
	defer logg.Sync()

	cfg := config.LoadConfig()
	server := app.BuildServer(logg, cfg)

	log.Println("Сервер запущен на 127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", server); err != nil {
		logg.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
