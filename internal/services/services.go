package services

import (
	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/stores"
)

type Services struct {
	UserService     *userService
	MaterialService *materialService
	PhraseService   *phraseService
	ChatService     *chatService
	MessageService  *messageService
	GeminiClient    *gemini.Client
}

func NewServices(s *stores.Stores, geminiClient *gemini.Client) *Services {
	return &Services{
		UserService:     &userService{store: s.UserStore},
		MaterialService: &materialService{store: s.MaterialStore},
		PhraseService:   &phraseService{store: s.PhraseStore, MaterialService: &materialService{store: s.MaterialStore}, GeminiClient: geminiClient},
		ChatService:     &chatService{chatStore: s.ChatStore, messageStore: s.MessageStore},
		MessageService:  &messageService{store: s.MessageStore, chatStore: s.ChatStore, geminiClient: geminiClient},
		GeminiClient:    geminiClient,
	}
}
