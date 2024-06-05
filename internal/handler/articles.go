package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

const (
	ErrInvalidArticleID       = "invalid article ID format"
	ErrInvalidArticleData     = "invalid article data"
	ErrForbiddenModify        = "forbidden to modify this article"
	ErrFailedUpdateArticle    = "failed to update article"
	ErrInvalidID              = "invalid ID"
	ErrFailedDeleteArticle    = "failed to delete article"
	ErrFailedRetrieveArticles = "failed to retrieve articles"
	ErrFailedCreateArticle    = "failed to create article"
	ErrArticleNotFound        = "article not found"
)

// ArticleHandler defines the methods for handling article-related requests.
type ArticleHandler interface {
	CreateArticle(c echo.Context) error
	GetArticleByID(c echo.Context) error
	UpdateArticle(c echo.Context) error
	DeleteArticle(c echo.Context) error
	GetAllArticles(c echo.Context) error
}

// articleHandler is the concrete implementation of ArticleHandler.
type articleHandler struct {
	services.ArticleService
}

func (h *articleHandler) GetArticleByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	article, err := h.ArticleService.GetArticleByID(uint(id), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}
	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := c.Bind(&article); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleData)
	}
	if err := validateArticle(&article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	article.UserID = userID
	log.Println(article.UserID)
	if err := h.ArticleService.CreateArticle(&article); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedCreateArticle)
	}
	return c.JSON(http.StatusCreated, article)
}

func (h *articleHandler) UpdateArticle(c echo.Context) error {
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}
	article, err := h.ArticleService.GetArticleByID(uint(id), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}
	if err := c.Bind(article); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleData)
	}
	if err := validateArticle(article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}
	if article.UserID != userID {
		return respondWithError(c, http.StatusForbidden, ErrForbiddenModify)
	}
	if err := h.ArticleService.UpdateArticle(uint(id), article); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedUpdateArticle)
	}
	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) DeleteArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidID)
	}
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	if err := h.ArticleService.DeleteArticle(uint(id), userID); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedDeleteArticle)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *articleHandler) GetAllArticles(c echo.Context) error {
	searchQuery := c.QueryParam("search")
	userID, err := getUserIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	articles, err := h.ArticleService.GetAllArticles(searchQuery, userID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedRetrieveArticles)
	}
	return c.JSON(http.StatusOK, articles)
}

func validateArticle(article *models.Article) error {
	validate := validator.New()
	if err := validate.Struct(article); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error in field '%s': %s", strings.ToLower(err.Field()), err.Tag())
			errorMessages = append(errorMessages, errorMessage)
		}
		return errors.New(strings.Join(errorMessages, ", "))
	}
	return nil
}
