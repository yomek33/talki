package models

import (
	"gorm.io/gorm"
)

const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

type Material struct {
	gorm.Model
	UserUID string   `gorm:"type:varchar(255);index;foreignKey" json:"uid"`
	Title   string   `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Content string   `gorm:"type:text" json:"content" `
	Phrases []Phrase `gorm:"foreignKey:MaterialID;references:ID"`
	Status  string   `gorm:"type:varchar(255)" json:"status"`
}
