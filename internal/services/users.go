package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userId uuid.UUID) error
}

type userService struct {
	store stores.UserStore
}


func (s *userService) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.store.CreateUser(user)
}


func (s *userService) GetUserByID(userId uuid.UUID) (*models.User, error) {
	return s.store.GetUserByID(userId)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.store.GetUserByEmail(email)
}

func (s *userService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	return s.store.UpdateUser(user)
}


func (s *userService) DeleteUser(userId uuid.UUID) error {

	return s.store.DeleteUser(userId)
}