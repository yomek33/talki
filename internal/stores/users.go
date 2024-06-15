package stores

import (
	"errors"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

const (
	ErrUserCannotBeNil = "user cannot be nil"
	ErrUserUIDRequired = "user ID required"
)

type UserStore interface {
	CreateUser(user *models.User) error
	GetUserByID(UserUID string) (*models.User, error)
	GetUserByUserUID(UserUID string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(UserUID string) error
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

func (store *userStore) GetUserByID(UserUID string) (*models.User, error) {
	var user models.User
	err := store.DB.Where("id = ?", UserUID).Preload("Materials").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (store *userStore) GetUserByUserUID(UserUID string) (*models.User, error) {
	var user models.User
	err := store.DB.Where("user_uid = ?", UserUID).First(&user).Error
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

func (store *userStore) DeleteUser(UserUID string) error {
	return store.PerformDBTransaction(func(tx *gorm.DB) error {
		// Delete the materials related to the user
		if err := tx.Where("user_uid = ?", UserUID).Delete(&models.Material{}).Error; err != nil {
			return err
		}
		// Delete the user
		if err := tx.Where("user_uid = ?", UserUID).Delete(&models.User{}).Error; err != nil {
			return err
		}
		return nil
	})
}
