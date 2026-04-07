package model

type User struct {
	ID       int  `json:"id" gorm:"primaryKey"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PhotoURL string `json:"photo_url"`
	Role     string `json:"role" gorm:"default:user"`
}

func (User) TableName() string {
	return "user"
}
