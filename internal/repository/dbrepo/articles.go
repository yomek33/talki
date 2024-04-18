package dbrepo

import (
	"errors"

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
	if article == nil {
		return errors.New("create article: article cannot be nil")
	}
	tx := repo.DB.Begin()
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *ArticleRepo) GetArticleByID(id uint, userID uint) (*models.Article, error) {
	var article models.Article
	err := repo.DB.Where("id = ? AND user_id = ?", id, userID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (repo *ArticleRepo) UpdateArticle(article *models.Article) error {
	if article == nil {
		return errors.New("update article: article cannot be nil")
	}
	return repo.DB.Save(article).Error
}

func (repo *ArticleRepo) DeleteArticle(id uint, userID uint) error {
	tx := repo.DB.Begin()
	if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Article{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (repo *ArticleRepo) GetAllArticles(searchQuery string, userID uint) ([]models.Article, error) {
	var articles []models.Article
	query := repo.DB.Where("user_id = ?", userID)
	if searchQuery != "" {
		query = query.Where("title LIKE ?", "%"+searchQuery+"%")
	}
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}
