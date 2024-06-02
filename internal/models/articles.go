package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	UserID  uuid.UUID `gorm:"index" json:"user_id"`
	Title   string    `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Content string    `gorm:"type:text" json:"content" `
}
