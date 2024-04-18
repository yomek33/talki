package repository

import (
	"github.com/yomek33/talki/internal/models"
)

type ArticleRepository interface {
	GetArticleByID(id uint, userID uint) (*models.Article, error)
	CreateArticle(article *models.Article) error
	UpdateArticle(id uint, article *models.Article) error
	DeleteArticle(id uint, userID uint) error
	GetAllArticles(searchQuery string, userID uint) ([]models.Article, error)
}
