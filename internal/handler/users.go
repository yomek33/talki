package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/models"
	"github.com/yomek33/talki/services"
)

const (
	ErrInvalidUserData = "invalid user data"
	ErrCouldNotCreateUser = "could not create user"
	ErrInvalidUserToken = "invalid user token"
	ErrUserNotFound = "user not found"
	ErrInvalidUserID = "invalid user ID"
	ErrCouldNotUpdateUser = "could not update user"
	ErrCouldNotDeleteUser = "could not delete user"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type userHandler struct {
	services.UserService
}


func (h *userHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidUserData})
	}
	if err := validateUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": ErrCouldNotCreateUser})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *userHandler) GetUserByID(c echo.Context) error {
	userID,err  := getUserIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": ErrUserNotFound})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	userID,err  := getUserIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidUserData})
	}
	if err := validateUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if user.ID != userID {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserID})
	}
	if err := h.UserService.UpdateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": ErrCouldNotUpdateUser})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	userID,err  := getUserIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}

	if err = h.UserService.DeleteUser(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": ErrCouldNotDeleteUser})
	}
	return c.NoContent(http.StatusNoContent)
}

func getUserIDByContext(c echo.Context) (uint, error) {
	id, ok := c.Get("user").(uint)
	if !ok {
		return 0, errors.New(ErrInvalidUserID)
	}
	return id, nil
}

func validateUser(user *models.User) error {
	validate := validator.New()
	errorMessages := make([]string, 0)
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error: %s", strings.ToLower(err.Field()))
			errorMessages = append(errorMessages, errorMessage)
		}
		if len(errorMessages) > 0 {
			return errors.New(strings.Join(errorMessages, ", "))
		}
	}
	return nil
}