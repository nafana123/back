package auth

import (
	"back/internal/config"
	authdto "back/internal/dto/auth"
	userdto "back/internal/dto/user"
	steamService "back/internal/service/steam"
	telegramService "back/internal/service/telegram"
	userService "back/internal/service/user"
	"back/internal/validator"
	"back/pkg/httputils"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type AuthHandler struct {
	logger          *zap.Logger
	cfg             *config.Config
	userService     userService.UserService
	telegramService telegramService.TelegramService
	steamService    steamService.SteamService
}

func NewAuthHandler(
	logger *zap.Logger,
	cfg *config.Config,
	userService userService.UserService,
	telegramService telegramService.TelegramService,
	steamService steamService.SteamService,
) *AuthHandler {
	return &AuthHandler{
		logger:          logger,
		cfg:             cfg,
		userService:     userService,
		telegramService: telegramService,
		steamService:    steamService,
	}
}

func (h *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	var data authdto.DataRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		h.logger.Error("Ошибка декодирования", zap.Error(err))
		http.Error(w, "Ошибка декодирования", http.StatusBadRequest)
		return
	}

	response, err := h.telegramService.TelegramAuth(data.Data, h.cfg.TelegramBotToken, h.cfg.JWTSecret)
	if err != nil {
		h.logger.Error("Ошибка авторизации пользователя", zap.Error(err))
		http.Error(w, "Ошибка авторизации пользователя", http.StatusBadRequest)
		return
	}

	httputils.RespondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) SteamAuth(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.steamService.GetAuthURL()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) SteamCallback(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	response, err := h.steamService.ValidateAndGetProfile(queryParams)

	if err != nil {
		h.logger.Error("Ошибка валидации Steam ответа", zap.Error(err))
		http.Error(w, "Validation failed", http.StatusInternalServerError)
		return
	}
	frontendURL := fmt.Sprintf("%s/auth/steam/callback?profile=%s", h.cfg.FrontendURL, response)
	http.Redirect(w, r, frontendURL, http.StatusFound)
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var req userdto.RegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msgError := userdto.ErrorResponse{Error: "Ошибка получения тела запроса"}
		h.logger.Error("Ошибка получения тела запроса", zap.Error(err))
		httputils.RespondJSON(w, http.StatusBadRequest, msgError)
		return
	}

	ok, msg, field := validator.ValidateRegistration(&req)
	if !ok {
		msgError := userdto.ErrorResponse{Error: msg, Field: field}
		h.logger.Error("Ошибка при валидации данных")
		httputils.RespondJSON(w, http.StatusBadRequest, msgError)
		return
	}

	token, err := h.userService.Register(req.Login, req.Email, req.Password)
	if err != nil {
		switch err {
		case userService.ErrLoginAlreadyExists:
			msgError := userdto.ErrorResponse{Error: "Логин уже занят", Field: "login"}
			httputils.RespondJSON(w, http.StatusConflict, msgError)
			return
		case userService.ErrEmailAlreadyExists:
			msgError := userdto.ErrorResponse{Error: "Email уже зарегистрирован", Field: "email"}
			httputils.RespondJSON(w, http.StatusConflict, msgError)
			return
		default:
			msgError := userdto.ErrorResponse{Error: "Внутренняя ошибка сервера", Field: "login"}
			h.logger.Error("Ошибка регистрации пользователя", zap.Error(err))
			httputils.RespondJSON(w, http.StatusInternalServerError, msgError)
			return
		}
	}

	httputils.RespondJSON(w, http.StatusCreated, token)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userdto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msgError := userdto.ErrorResponse{Error: "Ошибка получения тела запроса"}
		h.logger.Error("Ошибка получения тела запроса", zap.Error(err))
		httputils.RespondJSON(w, http.StatusBadRequest, msgError)
		return
	}

	ok, msg, field := validator.ValidateLogin(&req)
	if !ok {
		msgError := userdto.ErrorResponse{Error: msg, Field: field}
		h.logger.Error("Ошибка при валидации данных", zap.String("field", field))
		httputils.RespondJSON(w, http.StatusBadRequest, msgError)
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case userService.ErrInvalidCredentials:
			msgError := userdto.ErrorResponse{Error: "Неверный логин или пароль", Field: "login"}
			httputils.RespondJSON(w, http.StatusUnauthorized, msgError)
			return
		default:
			msgError := userdto.ErrorResponse{Error: "Внутренняя ошибка сервера"}
			h.logger.Error("Ошибка авторизации пользователя", zap.Error(err))
			httputils.RespondJSON(w, http.StatusInternalServerError, msgError)
			return
		}
	}

	httputils.RespondJSON(w, http.StatusCreated, token)
}