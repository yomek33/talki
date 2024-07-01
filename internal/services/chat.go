package services

import "github.com/yomek33/talki/internal/stores"

type ChatService interface {
	CreateChat() error
	GetChatByID() error
}

type chatService struct {
	store stores.ChatStore
}

type MessageService interface {
	CreateMessage() error
	GetMessages() error
}
