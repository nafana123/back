package main

import (
	"log"
	"net/http"

	"back/internal/config"
	"back/internal/database"
	"back/internal/logger"
	"back/internal/router"

	"go.uber.org/zap"
)

func main() {
	logg := logger.New()
	defer logg.Sync()

	cfg := config.LoadConfig()
	database.Connect(logg)
	server := router.New(logg, cfg)

	log.Println("Сервер запущен на 127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", server); err != nil {
		logg.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}