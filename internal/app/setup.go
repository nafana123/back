package app

import (
	"back/internal/cache"
	"back/internal/config"
	"back/internal/database"
	adminHandler "back/internal/handler/admin"
	credentialsHandler "back/internal/handler/auth/credentials"
	googleHandler "back/internal/handler/auth/google"
	steamHandler "back/internal/handler/auth/steam"
	telegramHandler "back/internal/handler/auth/telegram"
	gameHandler "back/internal/handler/game"
	mail "back/internal/mailer"
	gamerepo "back/internal/repository/game"
	participantrepo "back/internal/repository/participant"
	steamrepo "back/internal/repository/steam"
	tguserrepo "back/internal/repository/telegram"
	tournamentrepo "back/internal/repository/tournament"
	userrepo "back/internal/repository/user"
	router "back/internal/router"
	gameService "back/internal/service/game"
	googleService "back/internal/service/google"
	steamService "back/internal/service/steam"
	telegramService "back/internal/service/telegram"
	tournamentService "back/internal/service/tournament"
	userService "back/internal/service/user"
	"back/pkg/jwt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repositories struct {
	user        *userrepo.UserRepository
	telegram    *tguserrepo.TgUserRepository
	tournament  *tournamentrepo.TournamentRepository
	game        *gamerepo.GameRepository
	participant *participantrepo.ParticipantRepository
	steam       *steamrepo.SteamRepository
}

type services struct {
	user       *userService.Service
	telegram   *telegramService.Service
	steam      *steamService.SteamService
	tournament *tournamentService.TournamentService
	game       *gameService.GameService
	jwt        *jwt.Service
	google     *googleService.GoogleService
}

type handlers struct {
	credentials *credentialsHandler.Handler
	google      *googleHandler.Handler
	steam       *steamHandler.Handler
	telegram    *telegramHandler.Handler
	tournament  *adminHandler.TournamentHandler
	game        *gameHandler.Handler
}

func BuildServer(log *zap.Logger, cfg *config.Config) http.Handler {
	db := database.Connect(log, cfg)
	store := cache.NewMemoryStore(5 * time.Minute)
	mailer := mail.NewSMTPMailer(cfg)

	repositories := initRepositories(db)
	services := initServices(repositories, cfg, store, mailer)
	handlers := initHandlers(services, log, cfg)

	return router.New(router.RoutesDeps{
		CredentialsHandler: handlers.credentials,
		GoogleHandler:      handlers.google,
		SteamHandler:       handlers.steam,
		TelegramHandler:    handlers.telegram,
		JWTService:        services.jwt,
		TournamentHandler: handlers.tournament,
		GameHandler:       handlers.game,
	})
}

func initRepositories(db *gorm.DB) repositories {
	return repositories{
		user:        userrepo.NewUserRepository(db),
		telegram:    tguserrepo.NewTgUserRepository(db),
		tournament:  tournamentrepo.NewTournamentRepository(db),
		game:        gamerepo.NewGameRepository(db),
		participant: participantrepo.NewParticipantRepository(db),
		steam:       steamrepo.NewSteamRepository(db),
	}
}

func initServices(repositories repositories, cfg *config.Config, store *cache.MemoryStore, mailer *mail.SMTPMailer) services {
	return services{
		user:       userService.NewUserService(repositories.user, cfg.JWTSecret, store, mailer),
		telegram:   telegramService.NewTelegramService(repositories.telegram),
		steam:      steamService.NewSteamService(cfg, repositories.steam),
		tournament: tournamentService.NewTournamentService(repositories.tournament, repositories.participant),
		game:       gameService.NewGameService(repositories.game),
		jwt:        jwt.NewService(cfg.JWTSecret),
		google:     googleService.NewGoogleService(cfg, repositories.user),
	}
}

func initHandlers(services services, log *zap.Logger, cfg *config.Config) handlers {
	return handlers{
		credentials: credentialsHandler.NewHandler(log, services.user),
		google:      googleHandler.NewHandler(log, cfg, services.google),
		steam:       steamHandler.NewHandler(log, cfg, services.steam),
		telegram:    telegramHandler.NewHandler(log, cfg, services.telegram),
		tournament:  adminHandler.NewTournamentHandler(log, services.tournament),
		game:        gameHandler.NewHandler(log, services.game),
	}
}
