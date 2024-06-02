package services

import (
	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/stores"
)

type Services struct {
	UserService    *userService
	ArticleService *articleService
	PhraseService  *phraseService
	GeminiClient   *gemini.Client
}

func NewServices(s *stores.Stores, geminiClient *gemini.Client) *Services {
	return &Services{
		UserService:    &userService{store: s.UserStore},
		ArticleService: &articleService{store: s.ArticleStore},
		PhraseService:  &phraseService{store: s.PhraseStore},
		GeminiClient:   geminiClient,
	}
}

