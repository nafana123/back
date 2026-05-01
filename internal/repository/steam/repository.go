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
