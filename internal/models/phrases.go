package models

import (
	"gorm.io/gorm"
)

type Phrase struct {
	gorm.Model
	ID         int  `gorm:"primaryKey"`
	ArticleID  uint `gorm:"index"`
	Text       string
	Importance string
	Article    Article `gorm:"foreignKey:ArticleID"`
}
