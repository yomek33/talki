package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yomek33/talki/internal/services"
)

type Handlers struct {
	UserHandler
	ArticleHandler
}

func NewHandler(s *services.Services) *Handlers {
	return &Handlers{
		UserHandler:    &userHandler{UserService: s.UserService},
		ArticleHandler: &articleHandler{ArticleService: s.ArticleService},
	}
}

func (h *Handlers) SetDefault(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to our API")
	})
}

func (h *Handlers) SetAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")
	//middleware
	//api.Use(middleware.JWTMiddleware())

	api.POST("/articles", h.CreateArticle)
	api.GET("/articles", h.GetAllArticles)
	api.GET("/articles/:id", h.GetArticleByID)
	api.PUT("/articles/:id", h.UpdateArticle)
	api.DELETE("/articles/:id", h.DeleteArticle)

	api.POST("/users", h.CreateUser)
	api.GET("/users/:id", h.GetUserByID)
	api.PUT("/users/:id", h.UpdateUser)
	api.DELETE("/users/:id", h.DeleteUser)

}

func Echo() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	return e
}