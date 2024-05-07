package stores

import (
	"errors"

	"github.com/yomek33/talki/models"
	"gorm.io/gorm"
)

const (
	ErrUserCannotBeNil = "user cannot be nil"
	ErrUserIDRequired  = "user ID required"
)

type UserStore interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userStore struct {
	BaseStore
}


func (store *userStore) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("create user: " + ErrUserCannotBeNil)
	}

	return store.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}

func (store *userStore) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := store.DB.Where("user_id = ?", userID).Preload("Articles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (store *userStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := store.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (store *userStore) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("update user: " + ErrUserCannotBeNil)
	}

	return store.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Save(user).Error
	})
}

func (store *userStore) DeleteUser(userID uint) error {
	return store.PerformDBTransaction(func(tx *gorm.DB) error {
		// Soft delete
		return tx.Model(&models.User{}).Where("user_id = ?", userID).Update("deleted", true).Error
	})
}
