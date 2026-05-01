package router

import (
	admin "back/internal/handler/admin"
	auth "back/internal/handler/auth"
	game "back/internal/handler/game"
	"back/internal/handler/health"
	"back/internal/middleware"
	"back/pkg/jwt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type RoutesDeps struct {
	AuthHandler       *auth.AuthHandler
	JWTService        *jwt.Service
	TournamentHandler *admin.TournamentHandler
	GameHandler       *game.Handler
}

func New(deps RoutesDeps) http.Handler {
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
		r.Get("/health", health.HealthHandler)
		RegisterAuthRoutes(r, deps.AuthHandler, deps.JWTService)
		RegisterTournamentRoutes(r, deps.TournamentHandler, deps.JWTService)
		RegisterGameRoutes(r, deps.GameHandler)
	})

	return r
}

func RegisterAuthRoutes(r chi.Router, h *auth.AuthHandler, jwtService *jwt.Service) {
	r.Route("/auth", func(r chi.Router) {
		r.With(middleware.AuthRateLimit(10, time.Minute, "Слишком много попыток регистрации. Попробуйте позже.")).Post("/registration", h.Registration)
		r.With(middleware.AuthRateLimit(3, time.Minute, "Слишком много попыток подтверждения. Попробуйте позже.")).Post("/verify", h.Verify)
		r.With(middleware.AuthRateLimit(4, time.Minute, "Слишком много попыток входа. Попробуйте позже.")).Post("/login", h.Login)

		r.Post("/telegram", h.TelegramAuth)
		r.Get("/google", h.GoogleAuth)
		r.Get("/google/callback", h.GoogleCallback)
		r.Get("/steam/callback", h.SteamCallback)
	})

	r.With(middleware.Auth(jwtService)).Get("/auth/steam", h.SteamAuth)
}

func RegisterTournamentRoutes(r chi.Router, h *admin.TournamentHandler, jwtService *jwt.Service) {
	r.Get("/tournaments", h.GetTournaments)
	r.Get("/tournaments/{id}", h.GetTournament)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtService))

		r.Post("/tournaments/{id}/join", h.JoinTournament)
		r.Get("/tournaments/{id}/participants", h.GetParticipants)

		r.Route("/admin", func(r chi.Router) {
			r.Post("/tournaments", h.CreateTournament)
			r.Patch("/tournaments/{id}/status", h.ChangeTournamentStatus)
		})
	})
}

func RegisterGameRoutes(r chi.Router, h *game.Handler) {
	r.Get("/games", h.GetGames)
}
