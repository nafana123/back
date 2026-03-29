package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(log *zap.Logger) *http.Server {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}
}