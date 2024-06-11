package services

import (
	"errors"
	"fmt"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type ArticleService interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uint, UserUID string) (*models.Article, error)
	UpdateArticle(id uint, article *models.Article) error
	DeleteArticle(id uint, UserUID string) error
	GetAllArticles(searchQuery string, UserUID string) ([]models.Article, error)
}

type articleService struct {
	store stores.ArticleStore
}

var (
	ErrArticleNil          = errors.New("article cannot be nil")
	ErrMismatchedArticleID = errors.New("mismatched article ID")
)

func (s *articleService) CreateArticle(article *models.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	return s.store.CreateArticle(article)
}

func (s *articleService) GetArticleByID(id uint, UserUID string) (*models.Article, error) {
	article, err := s.store.GetArticleByID(id, UserUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get article by ID: %w", err)
	}
	return article, nil
}

func (s *articleService) UpdateArticle(id uint, article *models.Article) error {
	if article == nil {
		return ErrArticleNil
	}
	if id != article.ID {
		return ErrMismatchedArticleID
	}
	return s.store.UpdateArticle(id, article)
}

func (s *articleService) DeleteArticle(id uint, UserUID string) error {
	return s.store.DeleteArticle(id, UserUID)
}

func (s *articleService) GetAllArticles(searchQuery string, UserUID string) ([]models.Article, error) {
	return s.store.GetAllArticles(searchQuery, UserUID)
}
