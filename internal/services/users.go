package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByUserUID(UserUID string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(UserUID string) error
	CheckHashPassword(user *models.User, password string) bool
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

func (s *userService) GetUserByUserUID(UserUID string) (*models.User, error) {
	return s.store.GetUserByUserUID(UserUID)
}

func (s *userService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	return s.store.UpdateUser(user)
}

func (s *userService) DeleteUser(UserUID string) error {

	return s.store.DeleteUser(UserUID)
}

func (s *userService) CheckHashPassword(user *models.User, password string) bool {
	//err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return true
}

func PasswordEncrypt(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword), err
}
