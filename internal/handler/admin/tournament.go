package admin

import (
	"back/internal/dto"
	"back/internal/middleware"
	"back/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type TournamentHandler struct {
	logger            *zap.Logger
	tournamentService *service.TournamentService
}

func NewTournamentHandler(logger *zap.Logger, tournamentService *service.TournamentService) *TournamentHandler {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tournaments)
}

func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	tournament, err := h.tournamentService.GetTournament(tournamentID)

	if err != nil {
		h.logger.Error("Failed to get tournament", zap.Error(err))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tournament)
}

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req dto.CreateTournamentRequest

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	tournament, err := h.tournamentService.CreateTournament(&req)

	if err != nil {
		h.logger.Error("Failed to create tournament", zap.Error(err))
		return
	}

	response := dto.ToResponse(tournament)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TournamentHandler) ChangeTournamentStatus(w http.ResponseWriter, r *http.Request) {
	tournamentId := chi.URLParam(r, "id")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req struct {
		Status string `json:"status"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if err := h.tournamentService.ChangeStatus(tournamentId, req.Status); err != nil {
		http.Error(w, "Failed to change tournament status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  req.Status,
		"message": "success",
	})
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(participant)
}

func (h *TournamentHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	tournamentId := chi.URLParam(r, "id")

	participants, err := h.tournamentService.GetParticipantsByTournament(tournamentId)
	if err != nil {
		h.logger.Error("Failed to get participants", zap.Error(err))
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(participants)
}
