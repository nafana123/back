package repository

import "gorm.io/gorm"

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
