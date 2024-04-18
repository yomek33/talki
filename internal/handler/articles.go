package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/repository"
)

const (
    ErrInvalidArticleID = "Invalid article ID format"
    ErrInvalidArticleData = "Invalid article data"
    ErrForbiddenModify = "Forbidden to modify this article"
    ErrFailedUpdateArticle = "Failed to update article"
    ErrInvalidID = "Invalid ID"
    ErrFailedDeleteArticle = "Failed to delete article"
    ErrFailedRetrieveArticles = "Failed to retrieve articles"
	ErrFailedCreateArticle = "Failed to create article"
	ErrArticleNotFound = "Article not found"
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
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	article, err := h.ArticleRepo.GetArticleByID(uint(id), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}
	return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleData)
	}
	if err := h.ArticleRepo.CreateArticle(&article); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedCreateArticle)
	}
	return c.JSON(http.StatusCreated, article)
}


func (h *ArticleHandler) UpdateArticle(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
    }

    userID, err := getUserIDByContext(c)
    if err != nil {
        return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
    }
    var article models.Article
    if err := c.Bind(&article); err != nil {
        return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleData)
    }
    if article.UserID != userID {
        return respondWithError(c, http.StatusForbidden, ErrForbiddenModify)
    }
    if err := h.ArticleRepo.UpdateArticle(uint(id), &article); err != nil {
        return respondWithError(c, http.StatusInternalServerError, ErrFailedUpdateArticle)
    }
    return c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c echo.Context) error {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        return respondWithError(c, http.StatusBadRequest, ErrInvalidID)
    }
    userID, err := getUserIDByContext(c)
    if err != nil {
        return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
    }
    if err := h.ArticleRepo.DeleteArticle(uint(id), userID); err != nil {
        return respondWithError(c, http.StatusInternalServerError, ErrFailedDeleteArticle)
    }
    return c.NoContent(http.StatusNoContent)
}


func (h *ArticleHandler) GetAllArticles(c echo.Context) error {
    searchQuery := c.QueryParam("search")
    userID, err := getUserIDByContext(c)
    if err != nil {
        return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
    }
    articles, err := h.ArticleRepo.GetAllArticles(searchQuery, userID)
    if err != nil {
        return respondWithError(c, http.StatusInternalServerError, ErrFailedRetrieveArticles)
    }
    return c.JSON(http.StatusOK, articles)
}

func respondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}
