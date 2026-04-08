package router

import (
	"back/internal/handler/admin"
	auth "back/internal/handler/auth"
	"back/internal/handler/health"
	"back/internal/repository"
	"back/internal/service"
	"gorm.io/gorm"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func New(log *zap.Logger, authHandler *auth.AuthHandler, db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	tournamentRepo := repository.NewTournamentRepository(db)
	gameRepo := repository.NewGameRepository(db)

	tournamentService := service.NewTournamentService(tournamentRepo, gameRepo)
	gameService := service.NewGameService(tournamentRepo, gameRepo)

	tournamentHandler := admin.NewTournamentHandler(log, tournamentService)
	gameHandler := admin.NewGameHandler(log, gameService)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health.HealthHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/telegram", authHandler.TelegramAuth)
			r.Post("/registration", authHandler.Registration)
			r.Post("/login", authHandler.Login)
			r.Get("/steam", authHandler.SteamAuth)
			r.Get("/steam/callback", authHandler.SteamCallback)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Post("/tournaments", tournamentHandler.CreateTournament)
			r.Patch("/tournaments/{id}/status", tournamentHandler.ChangeTournamentStatus)

			r.Get("/games", gameHandler.GetGames) // todo сделать
		})

		r.Get("/tournaments", tournamentHandler.GetTournaments)
		r.Get("/tournaments/{id}/join", tournamentHandler.GetTournament)
		r.Get("/tournaments/{id}", tournamentHandler.GetTournament)
	})

	return r
}
