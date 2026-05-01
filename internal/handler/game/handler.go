package game

import (
	"back/internal/service/game"
	"back/pkg/httputils"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	logger      *zap.Logger
	gameService *game.GameService
}

func NewHandler(logger *zap.Logger, gameService *game.GameService) *Handler {
	return &Handler{
		logger:      logger,
		gameService: gameService,
	}
}

func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameService.GetAllGames()
	if err != nil {
		h.logger.Error("Ошибка получения списка игр", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить игры")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, games)
}
