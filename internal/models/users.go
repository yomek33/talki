package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   string `gorm:"type:varchar(255);primaryKey" json:"user_id"`
	Email    string `gorm:"type:varchar(255)" json:"email"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Password string `gorm:"type:varchar(255)" json:"password"`
	GoogleID string `gorm:"type:varchar(255)" json:"google_id"`
	Deleted  bool   `gorm:"default:false"`
	Articles []Article
}

type Dialogue struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	UserID       uint `gorm:"index"`
	InputText    string
	ResponseText string
	CreatedAt    time.Time
	User         User
}

type Progress struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	UserID       uint `gorm:"index"`
	PhraseID     uint `gorm:"index"`
	Status       string
	LastReviewed time.Time
	User         User
	Phrase       Phrase
}
