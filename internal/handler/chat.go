package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

type ChatHandler interface {
	CreateChat(c echo.Context) error
	GetChatByID(c echo.Context) error
}

type chatHandler struct {
	ChatService services.ChatService
}

func (h *chatHandler) CreateChat(c echo.Context) error {
	var chat models.Chat
	if err := c.Bind(chat); err != nil {
		logger.Errorf("Error binding chat data: %v", err)
		return respondWithError(c, http.StatusBadRequest, "Invalid chat data")
	}
	if err := h.ChatService.CreateChat(chat); err != nil {
		logger.Errorf("Error creating chat: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Could not create chat")
	}
	logger.Infof("Chat created successfully")
	return c.JSON(http.StatusCreated, chat)
}

func (h *chatHandler) GetChatByID(c echo.Context) error {
	id := c.Param("id")
	var chat models.Chat
	if err := h.ChatService.GetChatByID(id, &chat); err != nil {
		logger.Errorf("Chat room not found: %v", err)
		return respondWithError(c, http.StatusNotFound, "Chat room not found")
	}
	logger.Infof("Chat retrieved successfully")
	return c.JSON(http.StatusOK, chat)
}

// POST /chats/:id/messages
func (h *Handlers) CreateMessage(c echo.Context) error {
	chatID := c.Param("id")
	var message models.Message
	if err := c.Bind(message); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	message.chatID = chatID
	if err := h.DB.Create(&message).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, message)
}

// GET /chats/:id/messages
func (h *Handlers) GetMessages(c echo.Context) error {
	chatID := c.Param("id")
	var messages []models.Message
	if err := h.DB.Where("chat_room_id = ?", chatID).Find(&messages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, messages)
}
