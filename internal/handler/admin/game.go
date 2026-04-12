package admin

import (
	"back/internal/service/game"
	"back/pkg/httputils"
	"go.uber.org/zap"
	"net/http"
)

type GameHandler struct {
	Logger      *zap.Logger
	GameService *game.GameService
}

func NewGameHandler(logger *zap.Logger, gameService *game.GameService) *GameHandler {
	return &GameHandler{
		Logger:      logger,
		GameService: gameService,
	}
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.GameService.GetAllGames()

	if err != nil {
		h.Logger.Error("Failed to get all games", zap.Error(err))
		http.Error(w, "Failed to get games", http.StatusInternalServerError)
		return
	}

	httputils.RespondJSON(w, http.StatusOK, games)
}
