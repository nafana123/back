package tournament

import (
	"back/internal/model"

	"gorm.io/gorm"
)

type TournamentRepository struct {
	db *gorm.DB
}

func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	return &TournamentRepository{
		db: db,
	}
}

func (r *TournamentRepository) GetALl() ([]model.Tournament, error) {
	var tournaments []model.Tournament
	if err := r.db.Find(&tournaments).Error; err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (r *TournamentRepository) GetById(id string) (*model.Tournament, error) {
	var tournament model.Tournament

	if err := r.db.Where("id = ?", id).First(&tournament).Error; err != nil {
		return nil, err
	}

	return &tournament, nil
}

func (r *TournamentRepository) CreateTournament(tournament *model.Tournament) error {
	return r.db.Create(tournament).Error
}

func (r *TournamentRepository) Update(tournament *model.Tournament) error {
	return r.db.Save(tournament).Error
}
