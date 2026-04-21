package main

import (
	"back/internal/cache"
	"back/internal/handler/admin"
	"back/internal/service/game"
	tournamentService "back/internal/service/tournament"
	"back/pkg/jwt"
	"log"
	"net/http"
	"strings"
	"time"

	"back/internal/config"
	"back/internal/database"
	"back/internal/handler/auth"
	"back/internal/logger"
	mail "back/internal/mailer"
	"back/internal/repository"
	"back/internal/router"
	steamService "back/internal/service/steam"
	telegramService "back/internal/service/telegram"
	userService "back/internal/service/user"
	googleService "back/internal/service/google"

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
	if strings.TrimSpace(cfg.JWTSecret) == "" || cfg.JWTSecret == "default-secret-key" {
		logg.Fatal("JWT_SECRET не задан или небезопасен")
	}

	db := database.Connect(logg, cfg)

	store := cache.NewMemoryStore(5 * time.Minute)
	mailer := mail.NewSMTPMailer(cfg)

	userRepo := repository.NewUserRepository(db)
	tgUserRepo := repository.NewTgUserRepository(db)
	tournamentRepo := repository.NewTournamentRepository(db)
	gameRepo := repository.NewGameRepository(db)
	participantRepo := repository.NewParticipantRepository(db)

	userSvc := userService.NewUserService(userRepo, cfg.JWTSecret, store, mailer)
	telegramSvc:= telegramService.NewTelegramService(tgUserRepo)
	steamSvc := steamService.NewSteamService(cfg)
	tournamentSvc := tournamentService.NewTournamentService(tournamentRepo, participantRepo)
	gameSvc := game.NewGameService(gameRepo)
	jwtSvc:= jwt.NewService(cfg.JWTSecret)
	googleSvc := googleService.NewGoogleService(cfg, userRepo)

	authHandler := auth.NewAuthHandler(logg, cfg, userSvc, telegramSvc, steamSvc, googleSvc)
	tournamentHandler := admin.NewTournamentHandler(logg, tournamentSvc)
	gameHandler := admin.NewGameHandler(logg, gameSvc)

	server := router.New(logg, authHandler, jwtSvc, tournamentHandler, gameHandler)

	log.Println("Сервер запущен на 127.0.0.1:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", server); err != nil {
		logg.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}
