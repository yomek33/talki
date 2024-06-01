package stores

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

const (
	ErrArticleCannotBeNil = "article cannot be nil"
	ErrMismatchedArticleID = "article ID does not match the provided ID"
)

type ArticleStore interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uint, userID uuid.UUID) (*models.Article, error)
	UpdateArticle(id uint, article *models.Article) error
	DeleteArticle(id uint, userID uuid.UUID) error
	GetAllArticles(searchQuery string, userID uuid.UUID) ([]models.Article, error)
}

type articleStore struct {
	BaseStore
}

func (s *articleStore) CreateArticle(article *models.Article) error {
	if article == nil {
		return errors.New(ErrArticleCannotBeNil)
	}
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(article).Error
	})
}

func (s *articleStore) GetArticleByID(id uint, userID uuid.UUID) (*models.Article, error) {
	var article models.Article
	err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&article).Error
	return &article, err
}

func (s *articleStore) UpdateArticle(id uint, article *models.Article) error {
	if article == nil {
		return errors.New(ErrArticleCannotBeNil)
	}
	if id != article.ID {
		return errors.New(ErrMismatchedArticleID)
	}
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Model(&models.Article{}).Where("id = ?", id).Updates(article).Error
	})
}

func (s *articleStore) DeleteArticle(id uint, userID uuid.UUID) error {
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Article{}).Error
	})
}

func (s *articleStore) GetAllArticles(searchQuery string, userID uuid.UUID) ([]models.Article, error) {
	var articles []models.Article
	query := s.DB.Where("user_id = ?", userID)
	if searchQuery != "" {
		query = query.Where("title LIKE ?", "%"+searchQuery+"%")
	}
	err := query.Find(&articles).Error
	return articles, err
}