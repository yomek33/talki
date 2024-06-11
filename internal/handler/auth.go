package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/config"
	"github.com/yomek33/talki/internal/models"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

const serviceAccountJsonPath = "./service-account-credentials.json"

type Firebase struct {
	App        *firebase.App
	AuthClient *auth.Client
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
	return &Firebase{App: app, AuthClient: auth}, nil
}

func FirebaseAuthMiddleware(firebaseAuth *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/api/auth" {
				return next(c)
			}

			cookie, err := c.Cookie("session")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "No session cookie found")
			}

			decoded, err := firebaseAuth.VerifySessionCookieAndCheckRevoked(context.Background(), cookie.Value)
			if err != nil {
				if auth.IsSessionCookieRevoked(err) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Session cookie revoked")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session cookie")
			}

			// Pass the decoded token to the next handler
			c.Set("decodedToken", decoded)

			return next(c)
		}
	}
}

func (h *userHandler) GetGoogleLoginSignin(c echo.Context) error {
	setCORSHeaders(c)

	var req UserSignUpRequest
	if err := c.Bind(&req); err != nil {
		h.logError(err, "Invalid request payload")
		return respondWithError(c, http.StatusBadRequest, "Invalid request payload")
	}

	idToken := req.IDToken

	token, err := h.Firebase.AuthClient.VerifyIDToken(c.Request().Context(), idToken)
	if err != nil {
		h.logError(err, "Failed to verify ID token")
		return respondWithError(c, http.StatusUnauthorized, "Invalid ID token")
	}

	user, err := h.UserService.GetUserByGoogleID(token.UID)
	if err != nil && err != gorm.ErrRecordNotFound {
		h.logError(err, "GetUserByGoogleID error")
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		name, ok := token.Claims["name"].(string)
		if !ok {
			h.logError(err, "Name claim missing or invalid")
			return respondWithError(c, http.StatusBadRequest, "Invalid token claims")
		}
		user = &models.User{
			GoogleID: token.UID,
			Name:     name,
			UserID:   uuid.New(),
		}
		if err := h.UserService.CreateUser(user); err != nil {
			h.logError(err, "CreateUser error")
			return respondWithError(c, http.StatusInternalServerError, err.Error())
		}
		log.Printf("User created: %v", user)
	}

	cookieValue, err := h.Firebase.AuthClient.SessionCookie(c.Request().Context(), idToken, config.SessionDuration)
	if err != nil {
		h.logError(err, "Failed to create session cookie")
		return respondWithError(c, http.StatusInternalServerError, "Failed to create session cookie")
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Ensure Secure is set as per your environment
	}
	log.Println("Setting cookie:", cookie)
	c.SetCookie(cookie)

	log.Println("UserID: ", user.UserID)
	c.Set("user", user)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Success",
	})
}
