package services

import (
	"errors"
	"sync"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

// ChatService defines the interface for chat-related operations
type ChatService interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	GetChatByChatID(id uint, userUID string) (*models.Chat, error)
	UpdateChat(chat *models.Chat) error
	GetChatByMaterialID(materialID uint, userUID string) (*models.Chat, error)
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
func (s *chatService) GetChatByMaterialID(materialID uint, userUID string) (*models.Chat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.chatStore.GetChatByMaterialID(materialID, userUID)
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
