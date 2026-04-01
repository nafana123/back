package auth

import (
	"back/internal/config"
	authdto "back/internal/dto/auth"
	"back/pkg/httputils"
	"encoding/json"
	"net/http"

	authService "back/internal/service/auth"

	"go.uber.org/zap"
)

type AuthHandler struct {
	Logger *zap.Logger
	Cfg    *config.Config
}

func NewAuthHandler(logger *zap.Logger, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		Logger: logger,
		Cfg:    cfg,
	}
}

func (ah *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var data authdto.DataRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ah.Logger.Error("Ошибка декодирования", zap.Error(err))
		http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
		return
	}

	response, err := authService.TelegramAuth(data.Data, ah.Cfg.TelegramBotToken, ah.Cfg.TelegramBotSecret)
	if err != nil {
		ah.Logger.Error("Ошибка авторизации пользователя", zap.Error(err))
		http.Error(w, "Ошибка авторизации пользователя", http.StatusBadRequest)
		return
	}

	httputils.RespondJSON(w, http.StatusOK, response)
}