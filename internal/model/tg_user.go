package model

type TgUser struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
	Role         string `json:"role" gorm:"default:user"`
}

func (TgUser) TableName() string {
	return "tg_user"
}
