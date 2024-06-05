package handler

import (
	"context"
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
}

type phraseHandler struct {
	services.PhraseService
}

func (h *phraseHandler) GeneratePhrases(c echo.Context) error {
	// Parse articleID and userID from the request
	userID, err := getUserIDByContext(c)
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
	phrases, err := h.PhraseService.GeneratePhrases(ctx, uint(articleID), userID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, ErrFailedGeneratePhrases)
	}

	// Return the phrases as JSON response
	return c.JSON(http.StatusOK, phrases)
}
