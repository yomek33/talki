package stores

import (
	"errors"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type PhraseStore interface {
	CreatePhrase(phrase *models.Phrase) error
	GetPhrasesByMaterialID(materialID uint) ([]models.Phrase, error)
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
	if phrase.MaterialID == 0 {
		return errors.New("phrase MaterialID cannot be empty")
	}

	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(phrase).Error
	})
}
func (s *phraseStore) GetPhrasesByMaterialID(materialID uint) ([]models.Phrase, error) {
	var phrases []models.Phrase
	err := s.DB.Where("material_id = ?", materialID).Find(&phrases).Error
	return phrases, err
}
