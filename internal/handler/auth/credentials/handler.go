package credentials

import (
	authdto "back/internal/dto/auth"
	userdto "back/internal/dto/user"
	userService "back/internal/service/user"
	"back/internal/validator"
	"back/pkg/httputils"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	logger      *zap.Logger
	userService *userService.Service
}

func NewHandler(
	logger *zap.Logger,
	userService *userService.Service,
) *Handler {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.userService.Login(req.Email, req.Password)
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

	httputils.RespondJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.userService.CompleteVerification(req)
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

	httputils.RespondJSON(w, http.StatusCreated, resp)
}
