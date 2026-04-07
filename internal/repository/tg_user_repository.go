package repository

import (
	"back/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TgUserRepository interface {
	UpsertUser(user *model.TgUser) error
}

type tgUserRepository struct {
	db *gorm.DB
}

func NewTgUserRepository(db *gorm.DB) TgUserRepository {
	return &tgUserRepository{db: db}
}

func (r *tgUserRepository) UpsertUser(user *model.TgUser) error {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "username", "language_code", "photo_url"}),
	}).Create(user).Error
	return err
}