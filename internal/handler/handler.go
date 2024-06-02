package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yomek33/talki/internal/services"
)

type Handlers struct {
	UserHandler
	ArticleHandler
	jwtSecretKey string
}

func NewHandler(s *services.Services, jwtSecretKey string) *Handlers {
	return &Handlers{
		UserHandler:    &userHandler{UserService: s.UserService, jwtSecretKey: jwtSecretKey},
		ArticleHandler: &articleHandler{ArticleService: s.ArticleService},
		jwtSecretKey:   jwtSecretKey,
	}
}

func (h *Handlers) SetDefault(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to our API")
	})
}

func (h *Handlers) SetAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")

	// Public routes
	api.POST("/login", h.Login)
	api.POST("/users", h.CreateUser)

	// Protected routes
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

func Echo() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	return e
}

