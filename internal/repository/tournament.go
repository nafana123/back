package repository

import (
	"back/internal/model"
	"gorm.io/gorm"
)

type TournamentRepository struct {
	db *gorm.DB
}

func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	if db == nil {
		panic("db is nil")
	}

	return &TournamentRepository{
		db: db,
	}
}

func (r *TournamentRepository) GetALl() ([]model.Tournament, error) {
	var tournaments []model.Tournament

	result := r.db.Find(&tournaments)

	if result.Error != nil {
		return nil, result.Error
	}

	return tournaments, nil
}

func (r *TournamentRepository) GetById(id string) (*model.Tournament, error) {
	var tournament model.Tournament

	result := r.db.First(&tournament, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tournament, nil
}

func (r *TournamentRepository) CreateTournament(tournament *model.Tournament) error {
	return r.db.Create(tournament).Error
}

func (r *TournamentRepository) Update(tournament *model.Tournament) error {
	return r.db.Save(tournament).Error
}
