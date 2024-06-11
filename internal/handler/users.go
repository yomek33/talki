package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

const (
	ErrInvalidUserData     = "invalid user data"
	ErrCouldNotCreateUser  = "could not create user"
	ErrInvalidUserToken    = "invalid user token"
	ErrUserNotFound        = "user not found"
	ErrInvalidUserUID      = "invalid user ID"
	ErrCouldNotUpdateUser  = "could not update user"
	ErrCouldNotDeleteUser  = "could not delete user"
	ErrInvalidCredentials  = "invalid credentials"
	TokenExpirationMinutes = 60
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
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}
	if err := validateUser(&user); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotCreateUser)
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}
	if err := validateUser(&user); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}
	if user.UserUID != UserUID {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserUID)
	}
	if err := h.UserService.UpdateUser(&user); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotUpdateUser)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	if err = h.UserService.DeleteUser(UserUID); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotDeleteUser)
	}
	return c.NoContent(http.StatusNoContent)
}

func getUserUIDByContext(c echo.Context) (string, error) {
	user := c.Get("userUserUID").(string)
	if user == "" {
		return "", errors.New("invalid user token")
	}
	log.Println("UserUID: ", user)
	return user, nil
}

func validateUser(user *models.User) error {
	validate := validator.New()
	errorMessages := make([]string, 0)
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error in field '%s': %s", strings.ToLower(err.Field()), err.Tag())
			errorMessages = append(errorMessages, errorMessage)
		}
		if len(errorMessages) > 0 {
			return errors.New(strings.Join(errorMessages, ", "))
		}
	}
	return nil
}
