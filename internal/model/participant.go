package model

import "time"

type Participant struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	TournamentID string    `json:"tournament_id" gorm:"index"`
	UserID       int       `json:"user_id" gorm:"index"`
	JoinedAt     time.Time `json:"joined_at" gorm:"default:CURRENT_TIMESTAMP"`
	Place        int       `json:"place" gorm:"default:0"`
	PrizeWon     float64   `json:"prize_won" gorm:"default:0"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (Participant) TableName() string {
	return "participants"
}
