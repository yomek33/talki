package services

import (
	"context"
	"fmt"

	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type PhraseService interface {
	GeneratePhrases(ctx context.Context, articleID uint, UserUID string) ([]string, error)
}

type phraseService struct {
	store          stores.PhraseStore
	ArticleService *articleService
	GeminiClient   *gemini.Client
}

func (s *phraseService) StorePhrase(phrase *models.Phrase) error {
	return s.store.CreatePhrase(phrase)
}

func (s *phraseService) GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error) {
	return s.store.GetPhrasesByArticleID(articleID)
}

func GeneratePhrases(topic string) ([]string, error) {
	return []string{}, nil
}

func (s *phraseService) GeneratePhrases(ctx context.Context, articleID uint, UserUID string) ([]string, error) {
	article, err := s.ArticleService.GetArticleByID(articleID, UserUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch article: %w", err)
	}

	phrases, err := s.GeminiClient.GeneratePhrases(ctx, article.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate phrases: %w", err)
	}

	return phrases, nil
}

func (s *phraseService) StorePhrases(articleID uint, phrases []string) error {
	for _, phrase := range phrases {
		newPhrase := &models.Phrase{
			ArticleID:  articleID,
			Text:       phrase,
			Importance: determineImportance(phrase),
		}
		if err := s.store.CreatePhrase(newPhrase); err != nil {
			return fmt.Errorf("failed to store phrase: %w", err)
		}
	}

	return nil
}

func determineImportance(phrase string) string {
	return "high"
}
