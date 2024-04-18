package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Title     string	`gorm:"type:varchar(255)"`	
	Content   string    `gorm:"type:text"`
	User      User
}