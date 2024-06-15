package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/config"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

type UserSignUpRequest struct {
	IDToken     string `json:"id_token" validate:"required"`
	UID         string `json:"uid" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
}

func (h *userHandler) GetGoogleLoginSignin(c echo.Context) error {
	var req UserSignUpRequest
	if err := c.Bind(&req); err != nil {
		logger.Errorf("Error binding request: %v", err)
		return respondWithError(c, http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		logger.Errorf("Error validating request: %v \n Request: %v", err, req)
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.Firebase.AuthClient.VerifyIDToken(c.Request().Context(), req.IDToken)
	if err != nil {
		logger.Errorf("Error verifying ID token: %v", err)
		return respondWithError(c, http.StatusUnauthorized, "Invalid ID token")
	}

	user, err := h.UserService.GetUserByUserUID(token.UID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("Error getting user by UID: %v", err)
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		name, ok := token.Claims["name"].(string)
		if !ok {
			logger.Errorf("Invalid token claims")
			return respondWithError(c, http.StatusBadRequest, "Invalid token claims")
		}
		user = &models.User{
			UserUID: token.UID,
			Name:    name,
		}
		if err := h.UserService.CreateUser(user); err != nil {
			logger.Errorf("Error creating user: %v", err)
			return respondWithError(c, http.StatusInternalServerError, err.Error())
		}
	}

	sessionDuration := config.SessionDuration
	cookieValue, err := h.Firebase.AuthClient.SessionCookie(c.Request().Context(), req.IDToken, sessionDuration)
	if err != nil {
		logger.Errorf("Error creating session cookie: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Failed to create session cookie")
	}
	c.SetCookie(createSessionCookie(cookieValue))

	logger.Infof("Success, UserUID: %v", user.UserUID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}

func createSessionCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
}
