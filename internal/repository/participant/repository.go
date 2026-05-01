package participant

import (
	"back/internal/model"

	"gorm.io/gorm"
)

type ParticipantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{
		db: db,
	}
}

func (r *ParticipantRepository) Create(participant *model.Participant) error {
	return r.db.Create(participant).Error
}

func (r *ParticipantRepository) Exists(tournamentId string, userId int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Participant{}).Where("tournament_id = ? AND user_id = ?", tournamentId, userId).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ParticipantRepository) GetByTournament(tournamentId string) ([]model.Participant, error) {
	var participants []model.Participant

	err := r.db.Where("tournament_id = ?", tournamentId).Preload("User").Find(&participants).Error

	return participants, err
}
