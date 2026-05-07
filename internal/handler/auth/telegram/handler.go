package telegram

import (
	"back/internal/config"
	authdto "back/internal/dto/auth"
	telegramService "back/internal/service/telegram"
	"back/pkg/httputils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	logger          *zap.Logger
	cfg             *config.Config
	telegramService *telegramService.Service
}

func NewHandler(
	logger *zap.Logger,
	cfg *config.Config,
	telegramService *telegramService.Service,
) *Handler {
	return &Handler{
		logger:          logger,
		cfg:             cfg,
		telegramService: telegramService,
	}
}

func (h *Handler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	var req authdto.DataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		httputils.RespondDecodeError(w)
		return
	}

	token, err := h.telegramService.ValidateAuth(req.Data, h.cfg.TelegramBotToken, h.cfg.JWTSecret)
	if err != nil {
		h.logger.Error("Ошибка авторизации пользователя", zap.Error(err))
		httputils.RespondError(w, http.StatusBadRequest, "Ошибка авторизации пользователя")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, token)
}
