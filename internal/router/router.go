package router

import (
	admin "back/internal/handler/admin"
	credentials "back/internal/handler/auth/credentials"
	google "back/internal/handler/auth/google"
	steam "back/internal/handler/auth/steam"
	telegram "back/internal/handler/auth/telegram"
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
	CredentialsHandler *credentials.Handler
	GoogleHandler      *google.Handler
	SteamHandler       *steam.Handler
	TelegramHandler    *telegram.Handler
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
		RegisterAuthRoutes(r, deps.CredentialsHandler, deps.GoogleHandler, deps.SteamHandler, deps.TelegramHandler, deps.JWTService)
		RegisterTournamentRoutes(r, deps.TournamentHandler, deps.JWTService)
		RegisterGameRoutes(r, deps.GameHandler)
	})

	return r
}

func RegisterAuthRoutes(
	r chi.Router,
	credentialsHandler *credentials.Handler,
	googleHandler *google.Handler,
	steamHandler *steam.Handler,
	telegramHandler *telegram.Handler,
	jwtService *jwt.Service,
) {
	r.Route("/auth", func(r chi.Router) {
		r.With(middleware.AuthRateLimit(10, time.Minute, "Слишком много попыток регистрации. Попробуйте позже.")).Post("/registration", credentialsHandler.Registration)
		r.With(middleware.AuthRateLimit(3, time.Minute, "Слишком много попыток подтверждения. Попробуйте позже.")).Post("/verify", credentialsHandler.Verify)
		r.With(middleware.AuthRateLimit(4, time.Minute, "Слишком много попыток входа. Попробуйте позже.")).Post("/login", credentialsHandler.Login)

		r.Post("/telegram", telegramHandler.TelegramAuth)
		r.Get("/google", googleHandler.GoogleAuth)
		r.Get("/google/callback", googleHandler.GoogleCallback)
		r.Get("/steam/callback", steamHandler.SteamCallback)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtService))
		r.Post("/steam/logout", steamHandler.SteamLogout)
		r.Get("/steam/auth", steamHandler.SteamAuth)
	})
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
