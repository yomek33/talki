package models

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	MaterialID uint `gorm:"index"`
	Text       string
	Importance string
	Level      string
	Material   Material
}
