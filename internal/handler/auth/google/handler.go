package google

import (
	"back/internal/config"
	oauthstate "back/internal/oauthstate"
	googleService "back/internal/service/google"
	"back/pkg/httputils"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

type Handler struct {
	logger        *zap.Logger
	cfg           *config.Config
	googleService *googleService.GoogleService
}

func NewHandler(
	logger *zap.Logger,
	cfg *config.Config,
	googleService *googleService.GoogleService,
) *Handler {
	return &Handler{
		logger:        logger,
		cfg:           cfg,
		googleService: googleService,
	}
}

func (h *Handler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := oauthstate.Generate(h.cfg.JWTSecret)
	if err != nil {
		h.logger.Error("Ошибка генерации OAuth state", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось начать вход через Google")
		return
	}

	redirectURL := h.googleService.GetAuthURL(state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if err := oauthstate.Validate(state, h.cfg.JWTSecret); err != nil {
		h.logger.Error("Ошибка при валидации OAuth state", zap.Error(err))
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
