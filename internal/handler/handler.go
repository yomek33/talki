package handler

import (
	"fmt"
	"net/http"

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
		ArticleHandler: &articleHandler{ArticleService: s.ArticleService, PhraseService: s.PhraseService},
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

	r.POST("/articles", h.CreateArticle)
	r.GET("/articles", h.GetAllArticles)
	r.GET("/articles/:id", h.GetArticleByID)
	r.PUT("/articles/:id", h.UpdateArticle)
	r.DELETE("/articles/:id", h.DeleteArticle)
	r.POST("/articles", h.CreateArticle)
	r.GET("/articles/:id/status", h.CheckArticleStatus)
	r.GET("/articles/:id/phrases", h.GetProcessedPhrases)

	//r.GET("/users/:id", h.GetUserByID) TODO: GetUserByUserUID
	r.PUT("/users/:id", h.UpdateUser) //TODO: change to /users/:UserUID
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

	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
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
