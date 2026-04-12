package main

import (
	"back/internal/handler/admin"
	"back/internal/service/game"
	service2 "back/internal/service/tournament"
	"back/pkg/jwt"
	"log"
	"net/http"

	"back/internal/config"
	"back/internal/database"
	"back/internal/handler/auth"
	"back/internal/logger"
	"back/internal/repository"
	"back/internal/router"
	steamService "back/internal/service/steam"
	telegramService "back/internal/service/telegram"
	userService "back/internal/service/user"
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
	db := database.Connect(logg, cfg)

	userRepo := repository.NewUserRepository(db)
	tgUserRepo := repository.NewTgUserRepository(db)

	userSvc := userService.NewUserService(userRepo, cfg.JWTSecret)
	telegramSvc := telegramService.NewTelegramService(tgUserRepo)
	steamSvc := steamService.NewSteamService(cfg)

	authHandler := auth.NewAuthHandler(logg, cfg, userSvc, telegramSvc, steamSvc)
	jwtService := jwt.NewService(cfg.JWTSecret)

	tournamentRepo := repository.NewTournamentRepository(db)
	gameRepo := repository.NewGameRepository(db)
	participantRepo := repository.NewParticipantRepository(db)

	tournamentService := service2.NewTournamentService(tournamentRepo, participantRepo)
	gameService := game.NewGameService(gameRepo)

	tournamentHandler := admin.NewTournamentHandler(logg, tournamentService)
	gameHandler := admin.NewGameHandler(logg, gameService)

	server := router.New(logg, authHandler, db, jwtService, tournamentHandler, gameHandler)

	log.Println("Сервер запущен на 127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", server); err != nil {
		logg.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
