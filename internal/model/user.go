package model

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Login     string     `json:"login"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	PhotoURL  string     `json:"photo_url"`
	Role      string     `json:"role" gorm:"default:user"`
	SteamUser *SteamUser `json:"steam_user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string {
	return "user"
}
