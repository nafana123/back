package auth

import (
	"back/internal/config"
	authdto "back/internal/dto/auth"
	steamService "back/internal/service/auth/steam"
	telegramService "back/internal/service/auth/telegram"
	"back/pkg/httputils"
	"encoding/json"
	"fmt"
	"net/http"

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

func (ah *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	var data authdto.DataRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		ah.Logger.Error("Ошибка декодирования", zap.Error(err))
		http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
		return
	}

	response, err := telegramService.TelegramAuth(data.Data, ah.Cfg.TelegramBotToken, ah.Cfg.TelegramBotSecret)
	if err != nil {
		ah.Logger.Error("Ошибка авторизации пользователя", zap.Error(err))
		http.Error(w, "Ошибка авторизации пользователя", http.StatusBadRequest)
		return
	}

	httputils.RespondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) SteamAuth(w http.ResponseWriter, r *http.Request) {
	redirectURL := steamService.SteamParams()

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) SteamCallback(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	response, err := steamService.ValidateSteamResponse(queryParams)

	if err != nil {
		h.Logger.Error("Ошибка валидации Steam ответа", zap.Error(err))
		http.Error(w, "Validation failed", http.StatusInternalServerError)
		return
	}
	frontendURL := fmt.Sprintf("http://localhost:5173/auth/steam/callback?profile=%s", response)
	http.Redirect(w, r, frontendURL, http.StatusFound)
}
