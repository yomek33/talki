package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

const (
	ErrInvalidUserData     = "invalid user data"
	ErrCouldNotCreateUser  = "could not create user"
	ErrInvalidUserToken    = "invalid user token"
	ErrUserNotFound        = "user not found"
	ErrInvalidUserID       = "invalid user ID"
	ErrCouldNotUpdateUser  = "could not update user"
	ErrCouldNotDeleteUser  = "could not delete user"
	ErrInvalidCredentials  = "invalid credentials"
	TokenExpirationMinutes = 60
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	Login(c echo.Context) error
	GetGoogleLoginSignin(c echo.Context) error
}

type userHandler struct {
	services.UserService
	jwtSecretKey string
	Firebase     *Firebase
}

// JWT token
func (h *userHandler) generateJWTToken(userID uuid.UUID) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * TokenExpirationMinutes)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecretKey))
}

// Handlers
func (h *userHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}

	user.UserID = uuid.New()
	if err := validateUser(&user); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotCreateUser)
	}
	return c.JSON(http.StatusCreated, user)
}
func (h *userHandler) GetUserByID(c echo.Context) error {
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrUserNotFound)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	userID, err := getUserIDByContext(c)
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
	if user.UserID != userID {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserID)
	}
	if err := h.UserService.UpdateUser(&user); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotUpdateUser)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) DeleteUser(c echo.Context) error {
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	if err = h.UserService.DeleteUser(userID); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotDeleteUser)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *userHandler) Login(c echo.Context) error {
	var loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginRequest); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidUserData)
	}

	user, err := h.UserService.GetUserByEmail(loginRequest.Email)
	if err != nil || !h.UserService.CheckHashPassword(user, loginRequest.Password) {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidCredentials)
	}

	token, err := h.generateJWTToken(user.UserID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrCouldNotCreateUser)
	}

	// Set JWT token in HTTP-only cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Minute * TokenExpirationMinutes)
	cookie.HttpOnly = true
	cookie.Secure = false // Change to true in production
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{"message": "login successful"})
}

// Helper functions

func getUserIDByContext(c echo.Context) (uuid.UUID, error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
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
