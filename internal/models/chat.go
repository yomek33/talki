package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Detail         string    `gorm:"type:text" json:"detail"`
	MaterialID     uint      `gorm:"index" json:"material_id" validate:"required"`
	UserUID        string    `gorm:"index" json:"user_uid" validate:"required"`
	Messages       []Message `gorm:"foreignKey:ChatID;references:ID"`
	PendingMessage bool      `gorm:"default:false" json:"pending_message"`
}

type Message struct {
	gorm.Model
	ChatID     uint   `gorm:"index;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"chat_id" validate:"required"`
	Content    string `gorm:"type:text" json:"content"`
	UserUID    string `gorm:"index" json:"user_uid"`
	SenderType string `gorm:"type:varchar(255)" json:"sender_type"` // user or bot
}
