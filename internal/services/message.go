package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
	"gorm.io/gorm"
)

// MessageService defines the interface for message-related operations
type MessageService interface {
	CreateMessage(chatID uint, message *models.Message) (*models.Message, error)
	GetMessages(chatID uint) ([]models.Message, error)
	SendMessageToGemini(chatID uint, content, userUID string) (string, error)
}

type messageService struct {
	store        stores.MessageStore
	chatStore    stores.ChatStore
	geminiClient *gemini.Client
	mu           sync.Mutex
}

// NewMessageService creates a new instance of messageService
func NewMessageService(ms stores.MessageStore, cs stores.ChatStore, gc *gemini.Client) MessageService {
	return &messageService{
		store:        ms,
		chatStore:    cs,
		geminiClient: gc,
	}
}

func (s *messageService) CreateMessage(chatID uint, message *models.Message) (*models.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if message.Content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	message.ChatID = chatID
	return s.store.CreateMessage(message)
}

func (s *messageService) GetMessages(chatID uint) ([]models.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.store.GetMessages(chatID)
}

func (s *messageService) SendMessageToGemini(chatID uint, content, userUID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chat, err := s.chatStore.GetChatByChatID(chatID, userUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("chat not found")
		}
		return "", err
	}

	// if chat.PendingMessage {
	// 	return "", errors.New("previous message pending response")
	// }

	userMessage := &models.Message{
		ChatID:     chatID,
		UserUID:    userUID,
		Content:    content,
		SenderType: "user",
	}

	if _, err := s.store.CreateMessage(userMessage); err != nil {
    return "", err
	}
	logger.Infof("User message: %s", content)

	chat.PendingMessage ++
	if err := s.chatStore.UpdateChat(chat); err != nil {
		return "", err
	}

	logger.Infof("UpdateChat")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := s.geminiClient.SendMessageToGemini(ctx, chat, content)
	if err != nil {
		s.revertPendingMessageState(chat)
		return "", err
	}

	// response:="botbot"
	botMessage := &models.Message{
		ChatID:     chatID,
		Content:    response,
		SenderType: "bot",
	}

	if _, err := s.store.CreateMessage(botMessage); err != nil {
    return "", err
	}
	//chat.PendingMessage = false
	if err := s.chatStore.UpdateChat(chat); err != nil {
		return "", err
	}

	return response, nil
}

func (s *messageService) revertPendingMessageState(chat *models.Chat) {
	//chat.PendingMessage = false
	s.chatStore.UpdateChat(chat)
}
