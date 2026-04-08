package model

type Game struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`
	Icon        string `json:"icon"`
}

func (Game) TableName() string { return "game" }
