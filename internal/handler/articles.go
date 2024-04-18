package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/repository"
)

type ArticleHandler struct {
	ArticleRepo repository.ArticleRepository
}

func NewArticleHandler(repo repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{
		ArticleRepo: repo,
	}
}

func (h *ArticleHandler) HandleArticles(e *echo.Echo) {
	e.GET("/articles/:id", h.GetArticleByID)
	e.POST("/articles", h.CreateArticle)
	e.PUT("/articles/:id", h.UpdateArticle)
	e.DELETE("/articles/:id", h.DeleteArticle)
	e.GET("/articles", h.GetAllArticles)
}

func (h *ArticleHandler) GetArticleByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid article ID format")
	}
	userID, err := getUserIDbyToken(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}
	article, err := h.ArticleRepo.GetArticleByID(uint(id), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "Article not found")
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid article data")
	}
	if err := h.ArticleRepo.CreateArticle(&article); err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to create article")
	}
	return c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) UpdateArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid article ID format")
	}

	userID, err := getUserIDbyToken(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid article data")
	}
	if article.UserID != userID {
		return respondWithError(c, http.StatusForbidden, "Forbidden to modify this article")
	}
	if err := h.ArticleRepo.UpdateArticle(uint(id), &article); err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to update article")
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid ID")
	}
	userID, err := getUserIDbyToken(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}
	if err := h.ArticleRepo.DeleteArticle(uint(id), userID); err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to delete article")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ArticleHandler) GetAllArticles(c echo.Context) error {
	searchQuery := c.QueryParam("search")
	userID, err := getUserIDbyToken(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}
	articles, err := h.ArticleRepo.GetAllArticles(searchQuery, userID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to retrieve articles")
	}
	return c.JSON(http.StatusOK, articles)
}

func respondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}
