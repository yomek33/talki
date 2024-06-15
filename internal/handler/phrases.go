package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/services"
)

const (
	ErrFailedGeneratePhrases = "failed to generate phrases"
)

type PhraseHandler interface {
	GeneratePhrases(c echo.Context) error
	GetProcessedPhrases(c echo.Context) error
}

type phraseHandler struct {
	services.PhraseService
}

func (h *phraseHandler) GeneratePhrases(c echo.Context) error {
	log.Printf("Generating phrases for article")
	// Parse articleID and UserUID from the request
	UserUID, err := getUserUIDByContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	articleID, err := strconv.Atoi(c.Param("articleID"))
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call the service to generate phrases
	phrases, err := h.PhraseService.GeneratePhrases(ctx, uint(articleID), UserUID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedGeneratePhrases)
	}

	// Return the phrases as JSON response
	return c.JSON(http.StatusOK, phrases)
}

func (h *phraseHandler) GetProcessedPhrases(c echo.Context) error {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidArticleID)
	}

	phrases, err := h.PhraseService.GetPhrasesByArticleID(articleID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, phrases)
}
