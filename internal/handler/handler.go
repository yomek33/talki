package handler

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yomek33/talki/internal/services"
)

const frontendURI = "http://localhost:5173"

type Handlers struct {
	UserHandler
	ArticleHandler
	PhraseHandler
	jwtSecretKey string
	Firebase     *Firebase
}

func NewHandler(s *services.Services, jwtSecretKey string, firebase *Firebase) *Handlers {
	return &Handlers{
		UserHandler:    &userHandler{UserService: s.UserService, jwtSecretKey: jwtSecretKey, Firebase: firebase},
		ArticleHandler: &articleHandler{ArticleService: s.ArticleService},
		PhraseHandler:  &phraseHandler{PhraseService: s.PhraseService},
		jwtSecretKey:   jwtSecretKey,
		Firebase:       firebase,
	}
}

func (h *Handlers) SetDefault(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to our API")
	})
}

func (h *Handlers) SetAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")

	api.OPTIONS("/auth", handleOptions)
	api.POST("/auth", h.GetGoogleLoginSignin)

	r := api.Group("")
	r.Use(JWTMiddleware(h.jwtSecretKey))

	r.POST("/articles", h.CreateArticle)
	r.GET("/articles", h.GetAllArticles)
	r.GET("/articles/:id", h.GetArticleByID)
	r.PUT("/articles/:id", h.UpdateArticle)
	r.DELETE("/articles/:id", h.DeleteArticle)

	r.GET("/users/:id", h.GetUserByID)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)
}

func handleOptions(c echo.Context) error {
	setCORSHeaders(c)
	return c.NoContent(http.StatusNoContent)
}

func Echo() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{frontendURI},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-CSRF-Token",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieHTTPOnly: true,
	}))

	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
}

func JWTMiddleware(jwtSecretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Bypass middleware for OPTIONS requests
			if c.Request().Method == http.MethodOptions {
				return next(c)
			}

			// Extract token from cookie
			tokenString, err := c.Cookie("jwt")
			if err != nil || tokenString.Value == "" {
				return respondWithError(c, http.StatusUnauthorized, "missing or malformed jwt")
			}

			// Parse and validate token
			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
			})
			if err != nil || !token.Valid {
				return respondWithError(c, http.StatusUnauthorized, "invalid token")
			}

			return next(c)
		}
	}
}

func respondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var code = http.StatusInternalServerError
	var message interface{} = echo.Map{"message": "Internal Server Error"}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
		if he.Internal != nil {
			message = fmt.Sprintf("%v, %v", err, he.Internal)
		}
	}

	// Log the error
	c.Logger().Error(err)

	// Send JSON response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			c.NoContent(code)
		} else {
			c.JSON(code, message)
		}
	}
}

func setCORSHeaders(c echo.Context) {
	c.Response().Header().Set("Access-Control-Allow-Origin", frontendURI)
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
}

func (h *userHandler) logError(err error, msg string) {
	log.Printf("[ERROR] %s: %v\n", msg, err)
}
