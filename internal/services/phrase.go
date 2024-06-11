package services

import (
	"context"
	"fmt"
	"log"

	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type PhraseService interface {
	GeneratePhrases(ctx context.Context, articleID uint, UserUID string) ([]models.Phrase, error)
	StorePhrases(articleID uint, phrases []models.Phrase) error
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

func (s *phraseService) GeneratePhrases(ctx context.Context, articleID uint, UserUID string) ([]models.Phrase, error) {
	log.Println("Generating phrases")

	log.Println("ArticleID", articleID)
	log.Println("UserUID", UserUID)
	article, err := s.ArticleService.GetArticleByID(articleID, UserUID)
	log.Println("Article", article)
	if err != nil {
		log.Printf("Failed to fetch article: %v", err)
		return nil, fmt.Errorf("failed to fetch article: %w", err)
	}
	if article == nil {
		log.Printf("Article is nil")
		return nil, fmt.Errorf("article is nil")
	}

	// Check if GeminiClient is nil
	if s.GeminiClient == nil {
		log.Printf("GeminiClient is nil")
		return nil, fmt.Errorf("GeminiClient is nil")
	}

	log.Printf("Generating phrases for article %d", articleID)

	// Generate phrases using GeminiClient
	phraseTexts, err := s.GeminiClient.GeneratePhrases(ctx, article.Content)
	if err != nil {
		log.Printf("Failed to generate phrases: %v", err)
		return nil, fmt.Errorf("failed to generate phrases: %w", err)
	}
	if phraseTexts == nil {
		log.Printf("Generated phrases are nil")
		return nil, fmt.Errorf("generated phrases are nil")
	}

	var phrases []models.Phrase
	for _, phraseText := range phraseTexts {
		phrases = append(phrases, models.Phrase{
			ArticleID:  articleID,
			Text:       phraseText,
			Importance: determineImportance(phraseText),
			Article:    *article,
		})
	}

	// for i := range 10 {
	// 	phrases = append(phrases, models.Phrase{
	// 		ArticleID:  articleID,
	// 		Text:       fmt.Sprintf("phrase %d", i),
	// 		Importance: "high",
	// 	})
	// }
	return phrases, nil
}

func (s *phraseService) StorePhrases(articleID uint, phrases []models.Phrase) error {
	for _, phrase := range phrases {
		if err := s.store.CreatePhrase(&phrase); err != nil {
			return fmt.Errorf("failed to store phrase: %w", err)
		}
	}

	return nil
}

func determineImportance(phrase string) string {
	return "high"
}
