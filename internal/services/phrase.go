package services

import (
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type PhraseService interface {
	GeneratePhrases(topic string) ([]string, error)
}

type phraseService struct {
	store stores.PhraseStore
}

func (s *phraseService) StorePhrase(phrase *models.Phrase) error {
	return s.store.CreatePhrase(phrase)
}

func (s *phraseService) GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error) {
	return s.store.GetPhrasesByArticleID(articleID)
}
