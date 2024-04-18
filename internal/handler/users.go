package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/repository"
)

const (
	ErrInvalidUserData = "Invalid user data"
	ErrCouldNotCreateUser = "Could not create user"
	ErrInvalidUserToken = "Invalid user token"
	ErrUserNotFound = "User not found"
	ErrInvalidUserID = "Invalid user ID"
	ErrCouldNotUpdateUser = "Could not update user"
	ErrCouldNotDeleteUser = "Could not delete user"
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
	e.GET("/users/:id", h.GetUserByID)
	e.PUT("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": ErrInvalidUserData})
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
	if user.ID != userID {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": ErrInvalidUserID})
	}
	if err := h.UserRepo.UpdateUser(userID, &user); err != nil {
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