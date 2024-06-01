package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type ArticleService interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uint, userID uuid.UUID) (*models.Article, error)
	UpdateArticle(id uint, article *models.Article) error
	DeleteArticle(id uint, userID uuid.UUID) error
	GetAllArticles(searchQuery string, userID uuid.UUID) ([]models.Article, error)
}

type articleService struct {
	store stores.ArticleStore
}

func (s *articleService) CreateArticle(article *models.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	return s.store.CreateArticle(article)
}

func (s *articleService) GetArticleByID(id uint, userID uuid.UUID) (*models.Article, error) {
	return s.store.GetArticleByID(id, userID)
}

func (s *articleService) UpdateArticle(id uint, article *models.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	if id != article.ID {
		return errors.New("mismatched article ID")
	}
	return s.store.UpdateArticle(id, article)
}

func (s *articleService) DeleteArticle(id uint, userID uuid.UUID) error {
	return s.store.DeleteArticle(id, userID)
}

func (s *articleService) GetAllArticles(searchQuery string, userID uuid.UUID) ([]models.Article, error) {
	return s.store.GetAllArticles(searchQuery, userID)
}