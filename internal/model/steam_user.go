package model

type SteamUser struct {
	ID          int  `json:"id" gorm:"primaryKey"`
	UserID      int    `json:"user_id" gorm:"uniqueIndex;not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SteamID     string `json:"steam_id" gorm:"index;not null"`
	PersonaName string `json:"persona_name"`
	AvatarURL   string `json:"avatar_url"`
}

func (SteamUser) TableName() string {
	return "steam_user"
}
