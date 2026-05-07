package steam

import (
	"back/internal/model"

	"gorm.io/gorm"
)

type SteamRepository struct {
	db *gorm.DB
}

func NewSteamRepository(db *gorm.DB) *SteamRepository {
	return &SteamRepository{db: db}
}

func (r *SteamRepository) CreateSteamUser(steamUser *model.SteamUser) error {
	return r.db.Create(steamUser).Error
}

func (r *SteamRepository) DeleteSteamUser(userID int) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.SteamUser{}).Error
}
