package telegram

import (
	"back/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TgUserRepository struct {
	db *gorm.DB
}

func NewTgUserRepository(db *gorm.DB) *TgUserRepository {
	return &TgUserRepository{db: db}
}

func (r *TgUserRepository) UpsertUser(user *model.TgUser) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "username", "language_code", "photo_url"}),
	}).Create(user).Error
}
