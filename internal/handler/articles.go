package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

type ArticleHandler interface {
	CreateArticle(c echo.Context) error
	GetArticleByID(c echo.Context) error
	UpdateArticle(c echo.Context) error
	DeleteArticle(c echo.Context) error
	GetAllArticles(c echo.Context) error
	CheckArticleStatus(c echo.Context) error
}

type articleHandler struct {
	services.ArticleService
	services.PhraseService
}

func NewArticleHandler(articleService services.ArticleService, phraseService services.PhraseService) ArticleHandler {
	return &articleHandler{
		ArticleService: articleService,
		PhraseService:  phraseService,
	}
}

func (h *articleHandler) CreateArticle(c echo.Context) error {
	var article models.Article
	if err := bindAndValidate(c, &article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	article.UserUID = UserUID
	article.Status = "processing"

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id, err := h.ArticleService.CreateArticle(&article)
	if err != nil {
		logger.Errorf("Error creating article: %v, UserUID: %v", err, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedCreateArticle)
	}

	article.ID = id
	go h.processArticleAsync(ctx, article.ID, UserUID)

	logger.Info("Article created successfully")
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Article created successfully",
		"id":      article.ID,
	})
}

func (h *articleHandler) GetArticleByID(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}

	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	article, err := h.ArticleService.GetArticleByID(id, UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}

	logger.Infof("Retrieved article ArticleID;%v", id)
	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) UpdateArticle(c echo.Context) error {
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	articleID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}

	article, err := h.ArticleService.GetArticleByID(articleID, UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrArticleNotFound)
	}

	if err := bindAndValidate(c, article); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	if article.UserUID != UserUID {
		return respondWithError(c, http.StatusForbidden, ErrForbiddenModify)
	}

	if err := h.ArticleService.UpdateArticle(articleID, article); err != nil {
		logger.Errorf("Failed to update article: %v, ArticleID: %v, UserUID: %v", err, articleID, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedUpdateArticle)
	}

	logger.Infof("Updated article, ArticleID: %v, UserUID: %v", articleID, UserUID)
	return c.JSON(http.StatusOK, article)
}

func (h *articleHandler) DeleteArticle(c echo.Context) error {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidID)
	}

	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	if err := h.ArticleService.DeleteArticle(articleID, UserUID); err != nil {
		logger.Errorf("Failed to delete article: %v, ArticleID: %v, UserUID: %v", err, articleID, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedDeleteArticle)
	}

	logger.Infof("Deleted article, ArticleID: %v, UserUID: %v", articleID, UserUID)
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
		logger.Errorf("Failed to retrieve articles: %v, UserUID: %v", err, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedRetrieveArticles)
	}

	logger.Infof("Retrieved articles, ArticleCount: %v, UserUID: %v", len(articles), UserUID)
	return c.JSON(http.StatusOK, articles)
}

func (h *articleHandler) CheckArticleStatus(c echo.Context) error {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}

	status, err := h.ArticleService.GetArticleStatus(articleID)
	if err != nil {
		logger.Errorf("Failed to get article status: %v, ArticleID: %v", err, articleID)
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Infof("Checked article status, ArticleID: %v, Status: %v", articleID, status)
	return c.JSON(http.StatusOK, map[string]string{"status": status})
}

func (h *articleHandler) processArticleAsync(ctx context.Context, articleID uint, userUID string) {
	h.ArticleService.UpdateArticleStatus(articleID, "processing")

	phrases, err := h.PhraseService.GeneratePhrases(ctx, articleID, userUID)
	if err != nil {
		logger.Errorf("Failed to generate phrases: %v, ArticleID: %v, UserUID: %v", err, articleID, userUID)
		h.ArticleService.UpdateArticleStatus(articleID, "failed")
		return
	}

	if err = h.PhraseService.StorePhrases(articleID, phrases); err != nil {
		logger.Errorf("Failed to store phrases: %v, ArticleID: %v, UserUID: %v", err, articleID, userUID)
		h.ArticleService.UpdateArticleStatus(articleID, "failed")
		return
	}

	logger.Infof("Phrases generated and stored successfully, ArticleID: %v, UserUID: %v", articleID, userUID)
	h.ArticleService.UpdateArticleStatus(articleID, "completed")
}
