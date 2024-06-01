package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/models"
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

func (s *Services) GenerateAndStorePhrases(ctx context.Context, articleID uint, userID uuid.UUID) error {
	article, err := s.ArticleService.GetArticleByID(articleID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch article: %w", err)
	}

	phrases, err := s.GeminiClient.GeneratePhrases(ctx, article.Content)
	if err != nil {
		return fmt.Errorf("failed to generate phrases: %w", err)
	}

	for _, phrase := range phrases {
		newPhrase := &models.Phrase{
			ArticleID: article.ID,
			Text:      phrase,
			Importance: determineImportance(phrase),
		}
		if err := s.PhraseService.StorePhrase(newPhrase); err != nil {
			return fmt.Errorf("failed to store phrase: %w", err)
		}
	}

	return nil
}

func determineImportance(phrase string) string {
	return "high"
}
