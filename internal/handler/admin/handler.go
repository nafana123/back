package admin

import (
	tournamentdto "back/internal/dto/tournament"
	"back/internal/service/tournament"
	"back/pkg/httputils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	logger            *zap.Logger
	tournamentService *tournament.TournamentService
}

func NewHandler(logger *zap.Logger, tournamentService *tournament.TournamentService) *Handler {
	return &Handler{
		logger:            logger,
		tournamentService: tournamentService,
	}
}

func (h *Handler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req tournamentdto.CreateTournamentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.RespondDecodeError(w)
		return
	}

	t, err := h.tournamentService.CreateTournament(&req)

	if err != nil {
		h.logger.Error("Ошибка создания турнира", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось создать турнир")
		return
	}

	response := tournamentdto.ToResponse(t)

	httputils.RespondJSON(w, http.StatusCreated, response)
}

func (h *Handler) ChangeTournamentStatus(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	var req tournamentdto.StatusUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.RespondDecodeError(w)
		return
	}

	if err := h.tournamentService.ChangeStatus(tournamentID, req.Status); err != nil {
		h.logger.Error("Ошибка изменения статуса турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось изменить статус турнира")
		return
	}

	httputils.RespondNoContent(w)
}
