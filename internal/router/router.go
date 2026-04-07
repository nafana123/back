package router

import (
	"back/internal/handler/auth"
	"back/internal/handler/health"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func New(log *zap.Logger, authHandler *auth.AuthHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/health", health.HealthHandler)

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/telegram", authHandler.TelegramAuth)
		r.Post("/registration", authHandler.Registration)
		r.Post("/login", authHandler.Login)
		r.Get("/steam", authHandler.SteamAuth)
		r.Get("/steam/callback", authHandler.SteamCallback)
	})

	return r
}