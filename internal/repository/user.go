package repository

import (
	"back/internal/database"
	"back/internal/model"

	"gorm.io/gorm/clause"
)

func UpsertUser(user *model.User) error {
	err := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name","last_name","username","language_code","photo_url"}),
	}).Create(user).Error
	return err
}