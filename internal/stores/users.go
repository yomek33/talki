package stores

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

const (
	ErrUserCannotBeNil = "user cannot be nil"
	ErrUserIDRequired  = "user ID required"
)

type UserStore interface {
	CreateUser(user *models.User) error
	GetUserByID(userId uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByGoogleID(googleID string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userId uuid.UUID) error
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

func (store *userStore) GetUserByID(userID uuid.UUID) (*models.User, error) {
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
func (store *userStore) GetUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := store.DB.Where("google_id = ?", googleID).First(&user).Error
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

func (store *userStore) DeleteUser(userID uuid.UUID) error {
	return store.PerformDBTransaction(func(tx *gorm.DB) error {
		// Delete the articles related to the user
		if err := tx.Where("user_id = ?", userID).Delete(&models.Article{}).Error; err != nil {
			return err
		}
		// Delete the user
		if err := tx.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
			return err
		}
		return nil
	})
}
