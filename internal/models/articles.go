package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Title     string
	Content   string    `gorm:"type:text"`
	User      User
}