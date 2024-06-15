package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(255)" json:"name" validate:"required"`
	UserUID   string     `gorm:"type:varchar(255)" json:"user_uid" `
	Materials []Material `gorm:"foreignKey:UserUID;references:UserUID"`
}

type Dialogue struct {
	gorm.Model
	ID           int       `gorm:"primaryKey"`
	UserUID      string    `gorm:"index" validate:"required"`
	InputText    string    `validate:"requaired"`
	ResponseText string    `validate:"required"`
	CreatedAt    time.Time `validate:"required"`
}

type Progress struct {
	gorm.Model
	ID           int       `gorm:"primaryKey"`
	UserUID      string    `gorm:"index" validate:"required"`
	PhraseID     string    `gorm:"index" validate:"required"`
	Status       string    `validate:"required"`
	LastReviewed time.Time `validate:"required"`
	Phrase       Phrase
}
