package models

import (
	"gorm.io/gorm"
)

type Phrase struct {
	gorm.Model
	ID         int  `gorm:"primaryKey"`
	ArticleID  uint `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Text       string
	Importance string
}
