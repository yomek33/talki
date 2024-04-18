package repository

import (
	"github.com/yomek33/talki/internal/models"
)
type ArticleRepository interface {
	CreateArticle(article *models.Article) error
	GetArticleByID(id uint) (*models.Article, error)
	UpdateArticle(article *models.Article) error
	DeleteArticle(id uint) error
	ListArticles() ([]models.Article, error)
}