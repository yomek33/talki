package dbrepo

import (
	"errors"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

const (
	ErrUserCannotBeNil = "user cannot be nil"
	ErrUserIDRequired  = "user ID required"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (repo *UserRepo) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("create user: " + ErrUserCannotBeNil)
	}

	tx := repo.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserRepo) GetUserByID(userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("get user: " + ErrUserIDRequired)
	}
	var user models.User
	err := repo.DB.Where("user_id = ?", userID).Preload("Articles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepo) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("update user: " + ErrUserCannotBeNil)
	}

	tx := repo.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserRepo) DeleteUser(userID string) error {
	if userID == "" {
		return errors.New("delete user: " + ErrUserIDRequired)
	}

	tx := repo.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Soft delete
	if err := tx.Model(&models.User{}).Where("user_id = ?", userID).Update("deleted", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
