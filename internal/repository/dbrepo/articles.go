package dbrepo

import (
	"github.com/yomek33/talki/internal/models"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	DB *gorm.DB
}


func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		DB: db,
	}
}

func (repo *ArticleRepo) CreateArticle(article *models.Article) error {
	return repo.DB.Create(article).Error
}

func (repo *ArticleRepo) GetArticleByID(id uint) (*models.Article, error) {
	var article models.Article
	result := repo.DB.Preload("User").First(&article, id)
	return &article, result.Error
}

func (repo *ArticleRepo) UpdateArticle(article *models.Article) error {
	return repo.DB.Save(article).Error
}

func (repo *ArticleRepo) DeleteArticle(id uint) error {
	return repo.DB.Delete(&models.Article{}, id).Error
}

func (repo *ArticleRepo) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	result := repo.DB.Preload("User").Find(&articles)
	return articles, result.Error
}