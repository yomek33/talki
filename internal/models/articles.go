package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"id" validate:"required"`
	UserID  uuid.UUID   `gorm:"index" json:"user_id" validate:"required"`
	Title   string `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Content string `gorm:"type:text" json:"content" validate:"required"`
	User    User  `json:"user" validate:"required"`
}
