package handler

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
)

func FirebaseAuthMiddleware(firebaseAuth *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/api/auth" {
				return next(c)
			}

			cookie, err := c.Cookie("session")
			if err != nil {
				logger.Errorf("No session cookie found: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "No session cookie found")
			}

			decoded, err := firebaseAuth.VerifySessionCookieAndCheckRevoked(context.Background(), cookie.Value)
			if err != nil {
				if auth.IsSessionCookieRevoked(err) {
					logger.Errorf("Session cookie revoked: %v", err)
					return echo.NewHTTPError(http.StatusUnauthorized, "Session cookie revoked")
				}
				logger.Errorf("Invalid session cookie: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session cookie")
			}

			// Pass the decoded token to the next handler
			c.Set("decodedToken", decoded)
			c.Set("userUID", decoded.UID)
			logger.Infof("User authenticated, UserUID: %v", decoded.UID)
			return next(c)
		}
	}
}
