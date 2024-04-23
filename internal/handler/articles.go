package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/repository"
)

const (
    ErrInvalidArticleID = "invalid article ID format"
    ErrInvalidArticleData = "invalid article data"
    ErrForbiddenModify = "forbidden to modify this article"
    ErrFailedUpdateArticle = "failed to update article"
    ErrInvalidID = "invalid ID"
    ErrFailedDeleteArticle = "failed to delete article"
    ErrFailedRetrieveArticles = "failed to retrieve articles"
    ErrFailedCreateArticle = "failed to create article"
    ErrArticleNotFound = "article not found"
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
    if err := validateArticle(&article); err != nil {
        return respondWithError(c, http.StatusBadRequest, err.Error())
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
    if err := validateArticle(&article); err != nil {
        return respondWithError(c, http.StatusBadRequest, err.Error())
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


func validateArticle(article *models.Article) error {
	validate := validator.New()
	errorMessages := make([]string, 0)
	if err := validate.Struct(article); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error: %s", strings.ToLower(err.Field()))
			errorMessages = append(errorMessages, errorMessage)
		}
		if len(errorMessages) > 0 {
			return errors.New(strings.Join(errorMessages, ", "))
		}
	}
	return nil
}