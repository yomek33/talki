package models

import "gorm.io/gorm"

type Phrase struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey"`
	ArticleID  uint   `gorm:"index"`
	Text       string
	Importance string
	Article    Article
}