package admin

import (
	"back/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type GameHandler struct {
	Logger      *zap.Logger
	GameService *service.GameService
}

func NewGameHandler(logger *zap.Logger, gameService *service.GameService) *GameHandler {
	return &GameHandler{
		Logger:      logger,
		GameService: gameService,
	}
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {

}
