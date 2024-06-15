package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetGoogleLoginSignin(c echo.Context) error
}

type userHandler struct {
	services.UserService
	jwtSecretKey string
	Firebase     *Firebase
}

// Handlers
func (h *userHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		logger.Errorf("Error binding user data: %v", err)
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}
	if err := validateUser(&user); err != nil {
		logger.Errorf("Error validating user data: %v", err)
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		logger.Errorf("Error creating user: %v", err)
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotCreateUser)
	}
	logger.Infof("User created successfully")
	return c.JSON(http.StatusCreated, user)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		logger.Errorf("Error getting user UID from context: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		logger.Errorf("Error binding user data: %v", err)
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}
	if err := validateUser(&user); err != nil {
		logger.Errorf("Error validating user data: %v", err)
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}
	if user.UserUID != UserUID {
		logger.Errorf("Mismatch between user UID in context and request")
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserUID)
	}
	if err := h.UserService.UpdateUser(&user); err != nil {
		logger.Errorf("Error updating user: %v", err)
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotUpdateUser)
	}
	logger.Infof("User updated successfully")
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		logger.Errorf("Error getting user UID from context: %v", err)
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	if err = h.UserService.DeleteUser(UserUID); err != nil {
		logger.Errorf("Error deleting user: %v", err)
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotDeleteUser)
	}
	logger.Infof("User deleted successfully")
	return c.NoContent(http.StatusNoContent)
}
