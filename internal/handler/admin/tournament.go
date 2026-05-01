package admin

import (
	tournamentdto "back/internal/dto/tournament"
	"back/internal/middleware"
	"back/internal/service/tournament"
	"back/pkg/httputils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

func (h *TournamentHandler) respondDecodeError(w http.ResponseWriter, err error) {
	h.logger.Error("Ошибка декодирования тела запроса", zap.Error(err))
	httputils.RespondDecodeError(w)
}

func (h *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.tournamentService.GetAllTournaments()

	if err != nil {
		h.logger.Error("Ошибка получения турниров", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить турниры")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, tournaments)
}

func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	tournament, err := h.tournamentService.GetTournament(tournamentID)

	if err != nil {
		h.logger.Error("Ошибка получения турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить турнир")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, tournament)
}

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req tournamentdto.CreateTournamentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.respondDecodeError(w, err)
		return
	}

	tournament, err := h.tournamentService.CreateTournament(&req)

	if err != nil {
		h.logger.Error("Ошибка создания турнира", zap.Error(err))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось создать турнир")
		return
	}

	response := tournamentdto.ToResponse(tournament)

	httputils.RespondJSON(w, http.StatusCreated, response)
}

func (h *TournamentHandler) ChangeTournamentStatus(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	var req tournamentdto.StatusUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.respondDecodeError(w, err)
		return
	}

	if err := h.tournamentService.ChangeStatus(tournamentID, req.Status); err != nil {
		h.logger.Error("Ошибка изменения статуса турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось изменить статус турнира")
		return
	}

	httputils.RespondNoContent(w)
}

func (h *TournamentHandler) JoinTournament(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")
	userID, ok := middleware.GetUserID(r.Context())

	if !ok {
		httputils.RespondError(w, http.StatusUnauthorized, "Требуется авторизация")
		return
	}

	req := &tournamentdto.JoinTournamentRequest{
		TournamentID: tournamentID,
		UserID:       userID,
	}

	participant, err := h.tournamentService.JoinTournament(req)
	if err != nil {
		h.logger.Error("Ошибка присоединения к турниру", zap.Error(err), zap.String("tournament_id", tournamentID), zap.Int("user_id", userID))
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.RespondJSON(w, http.StatusCreated, participant)
}

func (h *TournamentHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	tournamentID := chi.URLParam(r, "id")

	participants, err := h.tournamentService.GetParticipantsByTournament(tournamentID)
	if err != nil {
		h.logger.Error("Ошибка получения участников турнира", zap.Error(err), zap.String("tournament_id", tournamentID))
		httputils.RespondError(w, http.StatusInternalServerError, "Не удалось получить участников турнира")
		return
	}

	httputils.RespondJSON(w, http.StatusOK, participants)
}
