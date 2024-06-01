package stores

import (
	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type PhraseStore interface {
	CreatePhrase(phrase *models.Phrase) error
	GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error)
}

type phraseStore struct {
	DB *gorm.DB
}

func (s *phraseStore) CreatePhrase(phrase *models.Phrase) error {
	return s.DB.Create(phrase).Error
}

func (s *phraseStore) GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error) {
	var phrases []models.Phrase
	err := s.DB.Where("article_id = ?", articleID).Find(&phrases).Error
	return phrases, err
}

func NewPhraseStore(db *gorm.DB) PhraseStore {
	return &phraseStore{DB: db}
}
