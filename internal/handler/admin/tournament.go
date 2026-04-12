package admin

import (
	"back/internal/dto"
	"back/internal/middleware"
	"back/internal/service/tournament"
	"back/pkg/httputils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type TournamentHandler struct {
	logger            *zap.Logger
	tournamentService *tournament.TournamentService
}

func NewTournamentHandler(logger *zap.Logger, tournamentService *tournament.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		logger:            logger,
		tournamentService: tournamentService,
	}
}

func (h *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.tournamentService.GetAllTournaments()

	if err != nil {
		h.logger.Error("Failed to get tournaments", zap.Error(err))

		return
	}

	httputils.RespondJSON(w, http.StatusOK, tournaments)
}

func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	tournament, err := h.tournamentService.GetTournament(tournamentID)

	if err != nil {
		h.logger.Error("Failed to get tournament", zap.Error(err))

		return
	}

	httputils.RespondJSON(w, http.StatusOK, tournament)
}

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateTournamentRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	tournament, err := h.tournamentService.CreateTournament(&body)

	if err != nil {
		h.logger.Error("Failed to create tournament", zap.Error(err))
		return
	}

	response := dto.ToResponse(tournament)

	httputils.RespondJSON(w, http.StatusCreated, response)
}

func (h *TournamentHandler) ChangeTournamentStatus(w http.ResponseWriter, r *http.Request) {
	tournamentId := chi.URLParam(r, "id")

	var body dto.StatusUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := h.tournamentService.ChangeStatus(tournamentId, body.Status); err != nil {
		http.Error(w, "Failed to change tournament status", http.StatusInternalServerError)
		return
	}

	httputils.RespondNoContent(w)
}

func (h *TournamentHandler) JoinTournament(w http.ResponseWriter, r *http.Request) {
	tournamentId := chi.URLParam(r, "id")
	userId, ok := middleware.GetUserID(r.Context())

	if !ok {
		http.Error(w, `{"error" : "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	req := &dto.JoinTournamentRequest{
		TournamentID: tournamentId,
		UserID:       userId,
	}

	participant, err := h.tournamentService.JoinTournament(req)
	if err != nil {
		h.logger.Error("Failed to join tournament", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httputils.RespondJSON(w, http.StatusCreated, participant)
}

func (h *TournamentHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	tournamentId := chi.URLParam(r, "id")

	participants, err := h.tournamentService.GetParticipantsByTournament(tournamentId)
	if err != nil {
		h.logger.Error("Failed to get participants", zap.Error(err))
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	httputils.RespondJSON(w, http.StatusCreated, participants)
}
