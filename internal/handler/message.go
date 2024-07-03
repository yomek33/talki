package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

type MessageHandler interface {
	CreateMessage(c echo.Context) error
	GetMessages(c echo.Context) error
}

type messageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(ms services.MessageService) MessageHandler {
	return &messageHandler{
		messageService: ms,
	}
}

// POST /chats/:chatId/message
func (h *messageHandler) CreateMessage(c echo.Context) error {
	chatID, err := parseUintParam(c, "chatId")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidChatID)
	}

	var message models.Message
	if err := c.Bind(&message); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid message data")
	}

	createdMessage, err := h.messageService.CreateMessage(chatID, &message)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to create message")
	}

	return c.JSON(http.StatusCreated, createdMessage)
}

// GET /chats/:chatId/messages
func (h *messageHandler) GetMessages(c echo.Context) error {
	chatID, err := parseUintParam(c, "chatId")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidChatID)
	}

	messages, err := h.messageService.GetMessages(chatID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "Failed to retrieve messages")
	}

	return c.JSON(http.StatusOK, messages)
}
