package router

import (
	"back/internal/config"
	auth "back/internal/handler/auth"
	"back/internal/handler/health"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func New(log *zap.Logger, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	authHandler := auth.NewAuthHandler(log, cfg)

	r.Get("/api/health", health.HealthHandler)
	r.Post("/api/auth/telegram", authHandler.Auth)

	r.Get("/api/auth/steam/login", authHandler.SteamLogin)
	r.Get("/api/auth/steam/callback", authHandler.SteamCallback)

	return r
}
