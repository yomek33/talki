package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/repository"
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
type UserHandler struct {
	UserRepo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepo: repo,
	}
}

func (h *UserHandler) HandleUsers(e *echo.Echo) {
	e.POST("/users", h.CreateUser)
	e.GET("/users", h.GetUserByID)
	e.PUT("/users", h.UpdateUser)
	e.DELETE("/users", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidUserData})
	}
	if err := validateUser(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.UserRepo.CreateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": ErrCouldNotCreateUser})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	userID,err  := getUserIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}
	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": ErrUserNotFound})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
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
	if err := h.UserRepo.UpdateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": ErrCouldNotUpdateUser})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID,err  := getUserIDByContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserToken})
	}

	if err = h.UserRepo.DeleteUser(userID); err != nil {
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