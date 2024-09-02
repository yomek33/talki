package services

import (
	"errors"
	"sync"

	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

// ChatService defines the interface for chat-related operations
type ChatService interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	GetChatByChatID(id uint, userUID string) (*models.Chat, error)
	UpdateChat(chat *models.Chat) error
	GetChatsByMaterialID(materialID uint, userUID string) ([]models.Chat, error)
}

// chatService implements the ChatService interface
type chatService struct {
	chatStore    stores.ChatStore
	messageStore stores.MessageStore
	mu           sync.RWMutex
}

// NewChatService creates a new instance of chatService
func NewChatService(cs stores.ChatStore, ms stores.MessageStore) ChatService {
	return &chatService{
		chatStore:    cs,
		messageStore: ms,
	}
}

// CreateChat creates a new chat
func (s *chatService) CreateChat(chat *models.Chat) (*models.Chat, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chatStore.CreateChat(chat)
}

// GetChatByMaterialID retrieves a chat by its material ID
func (s *chatService) GetChatsByMaterialID(materialID uint, userUID string) ([]models.Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	chats, err :=s.chatStore.GetChatsByMaterialID(materialID, userUID)
	if err != nil {
		logger.Info("error getting chats by material ID")
		newChat := models.Chat{
			MaterialID: materialID,
			UserUID:    userUID,
		}
		_, err := s.chatStore.CreateChat(&newChat)
		if err != nil {
			return nil, err
		}
		logger.Info("created new chat")
		chats = append(chats, newChat)
		return chats, nil
	}
	return chats, nil
}

// GetChatByChatID retrieves a chat by its ID
func (s *chatService) GetChatByChatID(id uint, userUID string) (*models.Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chat, err := s.chatStore.GetChatByChatID(id, userUID)
	if err != nil {
		return nil, err
	}

	if chat.UserUID != userUID {
		return nil, errors.New("unauthorized")
	}

	return chat, nil
}

// UpdateChat updates an existing chat
func (s *chatService) UpdateChat(chat *models.Chat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chatStore.UpdateChat(chat)
}
