package services

import (
	"errors"

	"github.com/yomek33/talki/models"
	"github.com/yomek33/talki/stores"
)

type ArticleService interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uint, userID uint) (*models.Article, error)
	UpdateArticle(id uint, article *models.Article) error
	DeleteArticle(id uint, userID uint) error
	GetAllArticles(searchQuery string, userID uint) ([]models.Article, error)
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

func (s *articleService) GetArticleByID(id uint, userID uint) (*models.Article, error) {
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

func (s *articleService) DeleteArticle(id uint, userID uint) error {
	return s.store.DeleteArticle(id, userID)
}

func (s *articleService) GetAllArticles(searchQuery string, userID uint) ([]models.Article, error) {
	return s.store.GetAllArticles(searchQuery, userID)
}