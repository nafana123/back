package model

import (
	"time"
)

type Tournament struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	Title          string    `json:"title"`
	PrizePool      float64   `json:"prize_pool"`
	EntryFee       float64   `json:"entry_fee"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status" gorm:"default:registration"` // registration, in_process, completed , failed
	Description    string    `json:"description"`
	Rules          string    `json:"rules"`
	MaxPlayers     uint      `json:"max_players"`
	CurrentPlayers uint      `json:"current_players"`
	CreatedAt      time.Time `json:"created_at"`
	Format         int       `json:"format"` // 1x1, 2x2, 3x3, 5x5 ...
	GameID         int       `json:"game_id" gorm:"index"`
}

const (
	StatusRegistration = "registration"
	StatusInProcess    = "in_process"
	StatusCompleted    = "completed"
	StatusFailed       = "failed"
)

func (Tournament) TableName() string { return "tournament" }
