package auth

import (
	"back/internal/config"
	authdto "back/internal/dto/auth"
	steamdto "back/internal/dto/steam"
	userdto "back/internal/dto/user"
	middleware "back/internal/middleware"
	oauthstate "back/internal/oauthstate"
	googleService "back/internal/service/google"
	steamService "back/internal/service/steam"
	telegramService "back/internal/service/telegram"
	userService "back/internal/service/user"
	"back/internal/validator"
	"back/pkg/httputils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

type AuthHandler struct {
	logger          *zap.Logger
	cfg             *config.Config
	userService     *userService.Service
	telegramService *telegramService.Service
	steamService    *steamService.SteamService
	googleService   *googleService.GoogleService
}

func NewAuthHandler(
	logger *zap.Logger,
	cfg *config.Config,
	userService *userService.Service,
	telegramService *telegramService.Service,
	steamService *steamService.SteamService,
	googleService *googleService.GoogleService,
) *AuthHandler {
	return &AuthHandler{
		logger:          logger,
		cfg:             cfg,
		userService:     userService,
		telegramService: telegramService,
		steamService:    steamService,
		googleService:   googleService,
	}
}

func (h *AuthHandler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) SteamAuth(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputils.RespondError(w, http.StatusUnauthorized, "Требуется авторизация")
		return
	}

	state, err := oauthstate.GenerateWithUserID(h.cfg.JWTSecret, userID)
	if err != nil {
		h.logger.Error("Ошибка генерации Steam state", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось начать вход через Steam")
		return
	}

	redirectURL := h.steamService.GetAuthURL(state)
	httputils.RespondJSON(w, http.StatusOK, steamdto.SteamRedirectResponse{RedirectURL: redirectURL})
}

func (h *AuthHandler) SteamCallback(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	state := queryParams.Get("state")
	userID, err := oauthstate.ValidateWithUserID(state, h.cfg.JWTSecret)
	if err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "Невалидный Steam state")
		return
	}

	err = h.steamService.ValidateCallback(queryParams, userID)
	if err != nil {
		h.logger.Error("Ошибка валидации Steam ответа", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Ошибка валидации ответа Steam")
		return
	}

	http.Redirect(w, r, h.cfg.FrontendURL, http.StatusFound)
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var req userdto.RegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		httputils.RespondDecodeError(w)
		return
	}

	ok, msg, field := validator.ValidateRegistration(&req)
	if !ok {
		h.logger.Error("Ошибка при валидации данных")
		httputils.RespondErrorWithField(w, http.StatusBadRequest, msg, field)
		return
	}

	err = h.userService.Register(req.Login, req.Email, req.Password)
	if err != nil {
		switch err {
		case userService.ErrLoginAlreadyExists:
			httputils.RespondErrorWithField(w, http.StatusConflict, "Логин уже занят", "login")
			return
		case userService.ErrEmailAlreadyExists:
			httputils.RespondErrorWithField(w, http.StatusConflict, "Email уже зарегистрирован", "email")
			return
		default:
			if errors.Is(err, userService.ErrEmailDeliveryFailed) {
				httputils.RespondErrorWithField(w, http.StatusBadRequest, "Не удалось отправить письмо: такой почты может не существовать или адрес указан с ошибкой", "email")
				return
			}
			h.logger.Error("Ошибка регистрации пользователя", zap.Error(err))
			httputils.RespondErrorWithField(w, http.StatusInternalServerError, "Внутренняя ошибка сервера", "login")
			return
		}
	}

	httputils.RespondNoContent(w)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req userdto.VerifyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		httputils.RespondDecodeError(w)
		return
	}

	ok, msg, field := validator.ValidateVerify(&req)
	if !ok {
		h.logger.Error("Ошибка при валидации данных", zap.String("field", field))
		httputils.RespondErrorWithField(w, http.StatusBadRequest, msg, field)
		return
	}

	token, err := h.userService.CompleteVerification(req)
	if err != nil {
		switch err {
		case userService.ErrInvalidVerificationCode:
			httputils.RespondErrorWithField(w, http.StatusConflict, "Неверный или просроченный код подтверждения", "code")
			return
		default:
			h.logger.Error("Ошибка подтверждения регистрации", zap.Error(err))
			httputils.RespondErrorWithField(w, http.StatusInternalServerError, "Внутренняя ошибка сервера", "code")
			return
		}
	}

	httputils.RespondJSON(w, http.StatusCreated, authdto.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userdto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
		httputils.RespondDecodeError(w)
		return
	}

	ok, msg, field := validator.ValidateLogin(&req)
	if !ok {
		h.logger.Error("Ошибка при валидации данных", zap.String("field", field))
		httputils.RespondErrorWithField(w, http.StatusBadRequest, msg, field)
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case userService.ErrGoogleOnlyAuth:
			httputils.RespondErrorWithField(w, http.StatusUnauthorized, "Этот аккаунт зарегистрирован через Google. Войдите с помощью Google.", "email")
			return
		case userService.ErrInvalidCredentials:
			httputils.RespondErrorWithField(w, http.StatusUnauthorized, "Неверный логин или пароль", "email")
			return
		default:
			h.logger.Error("Ошибка авторизации пользователя", zap.Error(err))
			httputils.RespondError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
			return
		}
	}

	httputils.RespondJSON(w, http.StatusCreated, token)
}

func (h *AuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := oauthstate.Generate(h.cfg.JWTSecret)
	if err != nil {
		h.logger.Error("Ошибка генерации OAuth state", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось начать вход через Google")
		return
	}

	redirectURL := h.googleService.GetAuthURL(state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if err := oauthstate.Validate(state, h.cfg.JWTSecret); err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "Невалидный OAuth state")
		return
	}

	response, err := h.googleService.GoogleValidate(r.Context(), code)
	if err != nil {
		h.logger.Error("Ошибка при валидации данных", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось выполнить вход через Google")
		return
	}

	baseFrontendURL := strings.TrimRight(h.cfg.FrontendURL, "/")
	frontendURL := fmt.Sprintf("%s/#token=%s", baseFrontendURL, url.QueryEscape(response.Token))
	http.Redirect(w, r, frontendURL, http.StatusFound)
}
