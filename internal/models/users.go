package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID   uuid.UUID `gorm:"type:varchar(255); primaryKey" json:"user_id" validate:"required"`
	Email    string `gorm:"type:varchar(255)" json:"email" validate:"required,email"`
	Name     string `gorm:"type:varchar(255)" json:"name" validate:"required"`
	Password string `gorm:"type:varchar(255)" json:"password" validate:"required"`
	GoogleID string `gorm:"type:varchar(255)" json:"google_id" `
	Deleted  bool   `gorm:"default:false"`
	Articles []Article
}

type Dialogue struct {
	gorm.Model
	ID           int `gorm:"primaryKey"`
	UserID       uuid.UUID `gorm:"index" validate:"required"`
	InputText    string `validate:"required"`
	ResponseText string `validate:"required"`
	CreatedAt    time.Time `validate:"required"`
	User         User `validate:"required"`
}

type Progress struct {
	gorm.Model
	ID          int `gorm:"primaryKey"`
	UserID       uuid.UUID `gorm:"index" validate:"required"`
	PhraseID     uuid.UUID `gorm:"index" validate:"required"`
	Status       string `validate:"required"`
	LastReviewed time.Time `validate:"required"`
	User         User `validate:"required"`
	Phrase       Phrase
}