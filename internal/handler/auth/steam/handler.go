package steam

import (
	"back/internal/config"
	steamdto "back/internal/dto/steam"
	middleware "back/internal/middleware"
	oauthstate "back/internal/oauthstate"
	steamService "back/internal/service/steam"
	"back/pkg/httputils"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	logger       *zap.Logger
	cfg          *config.Config
	steamService *steamService.SteamService
}

func NewHandler(
	logger *zap.Logger,
	cfg *config.Config,
	steamService *steamService.SteamService,
) *Handler {
	return &Handler{
		logger:       logger,
		cfg:          cfg,
		steamService: steamService,
	}
}

func (h *Handler) SteamAuth(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) SteamCallback(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) SteamLogout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputils.RespondError(w, http.StatusUnauthorized, "Требуется авторизация")
		return
	}

	err := h.steamService.Logout(userID)
	if err != nil {
		h.logger.Error("Ошибка выхода из Steam", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Ошибка выхода из Steam")
		return
	}

	httputils.RespondNoContent(w)
}
