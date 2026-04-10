package repository

import (
	"back/internal/model"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	if db == nil {
		panic("db is nil")
	}

	return &GameRepository{
		db: db,
	}
}

func (r *GameRepository) GetAll() ([]model.Game, error) {
	var games []model.Game

	result := r.db.Find(&games)

	if result.Error != nil {
		return nil, result.Error
	}

	return games, nil
}
