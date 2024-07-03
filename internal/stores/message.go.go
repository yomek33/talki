package stores

import (
	"errors"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type MessageStore interface {
	CreateMessage(message *models.Message) (*models.Message, error)
	GetMessages(chatID uint) ([]models.Message, error)
}

type messageStore struct {
	BaseStore
}

func (s *messageStore) CreateMessage(message *models.Message) (*models.Message, error) {
	if message == nil {
		return nil, errors.New("message cannot be nil")
	}

	// Check if the Chat exists
	var chat models.Chat
	if err := s.DB.First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chat does not exist")
		}
		return nil, err
	}

	// Create the message
	err := s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(message).Error
	})
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *messageStore) GetMessages(chatID uint) ([]models.Message, error) {
	var messages []models.Message
	err := s.DB.Where("chat_id = ?", chatID).Find(&messages).Error
	return messages, err
}
