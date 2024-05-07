package models

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	ArticleID  uint `gorm:"index"`
	Text       string
	Importance string
	Level      string
	Article    Article
}
