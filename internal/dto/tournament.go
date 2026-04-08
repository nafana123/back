package dto

import (
	"back/internal/model"
	"strconv"
	"time"
)

type TournamentBase struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	PrizePool      float64 `json:"prize_pool"`
	EntryFee       float64 `json:"entry_fee"`
	StartDate      string  `json:"start_date"`
	Status         string  `json:"status"`
	MaxPlayers     int     `json:"max_players"`
	CurrentPlayers int     `json:"current_players"`
	Format         string  `json:"format"`
	GameID         string  `json:"game_id"`
}

type CreateTournamentRequest struct {
	Title       string    `json:"title"`
	PrizePool   float64   `json:"prize_pool"`
	EntryFee    float64   `json:"entry_fee"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	MaxPlayers  int       `json:"max_players" validate:"min=2,max=10"`
	Format      int       `json:"format"`
	GameID      int       `json:"game_id"`
}

type TournamentResponse struct {
	TournamentBase
	EndDate     string    `json:"end_date"`
	Description string    `json:"description"`
	Rules       string    `json:"rules"`
	CreatedAt   time.Time `json:"created_at"`
}

type MinifiedTournamentResponse struct {
	TournamentBase
}

func (req *CreateTournamentRequest) ToModel() *model.Tournament {
	return &model.Tournament{
		Title:          req.Title,
		PrizePool:      req.PrizePool,
		EntryFee:       req.EntryFee,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Description:    req.Description,
		Rules:          req.Rules,
		MaxPlayers:     uint(req.MaxPlayers),
		Format:         req.Format,
		GameID:         req.GameID,
		CurrentPlayers: 0,
		Status:         "register",
		CreatedAt:      time.Now(),
	}
}

func ToResponse(t *model.Tournament) *TournamentResponse {
	return &TournamentResponse{
		TournamentBase: TournamentBase{
			ID:             t.ID,
			Title:          t.Title,
			PrizePool:      t.PrizePool,
			EntryFee:       t.EntryFee,
			StartDate:      t.StartDate.Format(time.RFC3339),
			Status:         t.Status,
			MaxPlayers:     int(t.MaxPlayers),
			CurrentPlayers: int(t.CurrentPlayers),
			Format:         strconv.Itoa(t.Format),
			GameID:         strconv.Itoa(t.GameID),
		},
		EndDate:     t.EndDate.Format(time.RFC3339),
		Description: t.Description,
		Rules:       t.Rules,
		CreatedAt:   t.CreatedAt,
	}
}
