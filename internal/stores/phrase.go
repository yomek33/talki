package stores

import (
	"errors"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type PhraseStore interface {
	CreatePhrase(phrase *models.Phrase) error
	GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error)
}

type phraseStore struct {
	BaseStore
}

func (s *phraseStore) CreatePhrase(phrase *models.Phrase) error {
	// バリデーション: phraseがnilでないことを確認
	if phrase == nil {
		return errors.New("phrase cannot be nil")
	}

	// バリデーション: 必要なフィールドが設定されていることを確認
	if phrase.Text == "" {
		return errors.New("phrase Text cannot be empty")
	}
	if phrase.ArticleID == 0 {
		return errors.New("phrase ArticleID cannot be empty")
	}

	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(phrase).Error
	})
}
func (s *phraseStore) GetPhrasesByArticleID(articleID uint) ([]models.Phrase, error) {
	var phrases []models.Phrase
	err := s.DB.Where("article_id = ?", articleID).Find(&phrases).Error
	return phrases, err
}
