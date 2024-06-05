package handler

import (
	"context"
	"fmt"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

const serviceAccountJsonPath = "./service-account-credentials.json"

type Firebase struct {
	App  *firebase.App
	Auth *auth.Client
}

type UserSignUpRequest struct {
	IDToken     string `json:"id_token"`
	UID         string `json:"uid"`
	DisplayName string `json:"display_name"`
}

func InitFirebase(ctx context.Context) (*Firebase, error) {
	opt := option.WithCredentialsFile(serviceAccountJsonPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}
	return &Firebase{App: app, Auth: auth}, nil
}

func (h *userHandler) exchangeCodeForToken(idToken string) (*auth.Token, error) {
	token, err := h.Firebase.Auth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}
	return token, nil
}

func (h *userHandler) GetGoogleLoginSignin(c echo.Context) error {
	setCORSHeaders(c)

	var req UserSignUpRequest
	if err := c.Bind(&req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request payload")
	}

	token, err := h.exchangeCodeForToken(req.IDToken)
	if err != nil {
		h.logError(err, "failed to verify ID token")
		return respondWithError(c, http.StatusUnauthorized, "Invalid ID token")
	}

	user, err := h.UserService.GetUserByGoogleID(token.UID)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.logError(err, "GetUserByGoogleID error")
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		user = &models.User{
			GoogleID: token.UID,
			Name:     token.Claims["name"].(string),
			UserID:   uuid.New(),
		}
		if err := h.UserService.CreateUser(user); err != nil {
			h.logError(err, "CreateUser error")
			return respondWithError(c, http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"jwt_token": req.IDToken})
}
