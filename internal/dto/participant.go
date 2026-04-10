package dto

import (
	"back/internal/model"
)

type JoinTournamentRequest struct {
	TournamentID string `json:"tournament_id"`
	UserID       int    `json:"user_id"`
}

func (req *JoinTournamentRequest) ToModel() *model.Participant {
	return &model.Participant{
		TournamentID: req.TournamentID,
		UserID:       req.UserID,
	}
}
