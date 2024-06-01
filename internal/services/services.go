package services

import (
	"context"

	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/stores"
)

// Services struct to hold different services
type Services struct {
	UserService    UserService
	ArticleService ArticleService
	GeminiClient   *gemini.Client
}

// NewServices initializes new services with the given stores and Gemini client
func NewServices(s *stores.Stores, geminiClient *gemini.Client) *Services {
	return &Services{
		UserService:    &userService{store: s.UserStore},
		ArticleService: &articleService{store: s.ArticleStore},
		GeminiClient:   geminiClient,
	}
}

// Example function using Gemini client
func (s *Services) GeneratePhrases(ctx context.Context, topic string) ([]string, error) {
	return s.GeminiClient.GeneratePhrases(ctx, topic)
}
