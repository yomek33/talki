package repository

import "github.com/yomek33/talki/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uint, user *models.User) error
	DeleteUser(id uint) error
}
