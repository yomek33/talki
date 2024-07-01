package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Detail     string    `gorm:"type:text" json:"detail" `
	MaterialID uint      `gorm:"index" json:"material_id" validate:"required"`
	UserID     uint      `gorm:"index" json:"user_id" validate:"required"`
	Messages   []Message `gorm:"foreignKey:ChatRoomID;references:ID"`
}

type Message struct {
	gorm.Model
	ChatID     uint   `gorm:"index" json:"chat_id" validate:"required"`
	Content    string `gorm:"type:text" json:"content"`
	UserID     uint   `gorm:"index" json:"sender_id"`
	SenderType string `gorm:"type:varchar(255)" json:"sender_type"` // user or bot
}
