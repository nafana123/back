package tournament

import (
	"back/internal/middleware"
	tournamentsvc "back/internal/service/tournament"
	"back/pkg/httputils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	logger            *zap.Logger
	tournamentService *tournamentsvc.TournamentService
}

func NewHandler(logger *zap.Logger, tournamentService *tournamentsvc.TournamentService) *Handler {
	return &Handler{
		logger:            logger,
		tournamentService: tournamentService,
	}
}

func (h *Handler) GetTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.tournamentService.GetAllTournaments()

	if err != nil {
		h.logger.Error("Ошибка получения турниров", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить турниры")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, tournaments)
}

func (h *Handler) GetTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	t, err := h.tournamentService.GetTournament(tournamentID)

	if err != nil {
		h.logger.Error("Ошибка получения турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить турнир")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, t)
}

func (h *Handler) JoinTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputils.RespondError(w, http.StatusUnauthorized, "Требуется авторизация")
		return
	}

	participant, err := h.tournamentService.JoinTournament(tournamentID, userID)
	if err != nil {
		h.logger.Error("Ошибка присоединения к турниру", zap.Error(err), zap.String("tournament_id", tournamentID), zap.Int("user_id", userID))
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.RespondJSON(w, http.StatusCreated, participant)
}

func (h *Handler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	participants, err := h.tournamentService.GetParticipantsByTournament(tournamentID)
	if err != nil {
		h.logger.Error("Ошибка получения участников турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить участников турнира")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, participants)
}
