package main

import (
	"log"

	"back/internal/logger"
	"back/internal/router"
)

func main() {
	logg := logger.New()
	defer logg.Sync()

	server := router.New(logg)

	log.Println("Сервер запущен", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}