package stores

import (
	"errors"
	"log"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type ChatStore interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	GetChatByChatID(id uint, UserUID string) (*models.Chat, error)
	GetChatByMaterialID(materialID uint, userUID string) (*models.Chat, error)
	UpdateChat(chat *models.Chat) error
}

type chatStore struct {
	BaseStore
}

func (s *chatStore) CreateChat(chat *models.Chat) (*models.Chat, error) {
	if chat == nil {
		return nil, errors.New("chat cannot be nil")
	}
	err := s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(chat).Error
	})
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *chatStore) GetChatByChatID(id uint, UserUID string) (*models.Chat, error) {
	log.Println("store chat id", id)
	var chat models.Chat
	err := s.DB.Where("id = ? AND user_uid = ?", id, UserUID).Preload("Messages").First(&chat).Error
	return &chat, err
}

func (s *chatStore) GetChatByMaterialID(materialID uint, userUID string) (*models.Chat, error) {
	var chat models.Chat
	err := s.DB.Where("material_id = ? AND user_uid = ?", materialID, userUID).Preload("Messages").First(&chat).Error
	return &chat, err
}

func (r *chatStore) UpdateChat(chat *models.Chat) error {
	if chat == nil {
		return errors.New("chat cannot be nil")
	}
	return r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.Chat{}).Where("id = ?", chat.ID).Updates(chat).Error
	})
}
