package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	services.PhraseService
}

func (h *articleHandler) GetArticleByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	article, err := h.ArticleService.GetArticleByID(uint(id), UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}
	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := bindAndValidateArticle(c, &article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	article.UserUID = UserUID

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	resultChan := make(chan error, 1)
	go func() {
		id, err := h.ArticleService.CreateArticle(&article)
		if err != nil {
			log.Printf("Error creating article: %v\n", err)
			resultChan <- err
			return
		}
		article.ID = id
		resultChan <- nil
	}()

	select {
	case <-ctx.Done():
		return respondWithError(c, http.StatusGatewayTimeout, "request timed out while creating article")
	case err := <-resultChan:
		if err != nil {
			return respondWithError(c, http.StatusInternalServerError, ErrFailedCreateArticle)
		}
	}

	var phrases []models.Phrase
	phraseErrChan := make(chan error, 1)
	phrasesChan := make(chan []models.Phrase, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\n", r)
				phraseErrChan <- fmt.Errorf("panic occurred: %v", r)
			}
		}()

		phrases, err = h.PhraseService.GeneratePhrases(ctx, article.ID, UserUID)
		if err != nil {
			log.Printf("Failed to generate phrases: %v\n", err)
			phraseErrChan <- err
			return
		}
		if phrases == nil {
			log.Printf("Generated phrases are nil")
			phraseErrChan <- fmt.Errorf("generated phrases are nil")
			return
		}

		err = h.PhraseService.StorePhrases(article.ID, phrases)
		if err != nil {
			log.Printf("Failed to store phrases: %v\n", err)
			phraseErrChan <- err
			return
		}
		phraseErrChan <- nil
		phrasesChan <- phrases
	}()

	log.Println("Waiting for phrases to be generated and stored")
	log.Println("send phrases to the client", phrases)

	phrase := <-phrasesChan
	if err := <-phraseErrChan; err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedGeneratePhrases)
	}
	return c.JSON(http.StatusCreated, phrase)
}

func bindAndValidateArticle(c echo.Context, article *models.Article) error {
	if err := c.Bind(article); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidArticleData, err)
	}
	if err := validateArticle(article); err != nil {
		return err
	}
	return nil
}

func (h *articleHandler) UpdateArticle(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	articleID := c.Param("id")
	id, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}
	article, err := h.ArticleService.GetArticleByID(uint(id), UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}
	if err := c.Bind(article); err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleData)
	}
	if err := validateArticle(article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}
	if article.UserUID != UserUID {
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
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	if err := h.ArticleService.DeleteArticle(uint(id), UserUID); err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedDeleteArticle)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *articleHandler) GetAllArticles(c echo.Context) error {
	searchQuery := c.QueryParam("search")
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}
	articles, err := h.ArticleService.GetAllArticles(searchQuery, UserUID)
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
