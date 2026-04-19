package router

import (
	admin "back/internal/handler/admin"
	auth "back/internal/handler/auth"
	"back/internal/handler/health"
	"back/internal/middleware"
	"back/pkg/jwt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func New(log *zap.Logger, authHandler *auth.AuthHandler, jwtService *jwt.Service, tournamentHandler *admin.TournamentHandler, gameHandler *admin.GameHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		// Публичные роуты
		r.Get("/health", health.HealthHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/telegram", authHandler.TelegramAuth)
			r.Post("/registration", authHandler.Registration)
			r.Post("/verify", authHandler.Verify)
			r.Post("/login", authHandler.Login)
			r.Get("/steam", authHandler.SteamAuth)
			r.Get("/steam/callback", authHandler.SteamCallback)
		})

		r.Get("/tournaments", tournamentHandler.GetTournaments)
		r.Get("/tournaments/{id}", tournamentHandler.GetTournament)
		r.Get("/games", gameHandler.GetGames)

		// Роуты требующие авторизацию
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtService))

			r.Post("/tournaments/{id}/join", tournamentHandler.JoinTournament)
			r.Get("/tournaments/{id}/participants", tournamentHandler.GetParticipants)

			r.Route("/admin", func(r chi.Router) {
				r.Post("/tournaments", tournamentHandler.CreateTournament)
				r.Patch("/tournaments/{id}/status", tournamentHandler.ChangeTournamentStatus)
			})
		})
	})

	return r
}
