package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

// ChatHandler defines the interface for chat-related operations
type ChatHandler interface {
	CreateChat(c echo.Context) error
	GetChatByMaterialID(c echo.Context) error
	GetChatByChatID(c echo.Context) error
	ChatWithGemini(c echo.Context) error
}

// chatHandler implements the ChatHandler interface
type chatHandler struct {
	chatService     services.ChatService
	messageService  services.MessageService
	materialService services.MaterialService
}

// NewChatHandler creates a new instance of chatHandler
func NewChatHandler(cs services.ChatService, ms services.MessageService) ChatHandler {
	return &chatHandler{
		chatService:    cs,
		messageService: ms,
	}
}

// CreateChat handles the creation of a new chat
func (h *chatHandler) CreateChat(c echo.Context) error {
	userUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}

	var chat models.Chat
	if err := c.Bind(&chat); err != nil {
		logger.Errorf("Error binding chat data: %v", err)
		return respondWithError(c, http.StatusBadRequest, "Invalid chat data")
	}

	chat.UserUID = userUID
	chat.CreatedAt = time.Now()

	createdChat, err := h.chatService.CreateChat(&chat)
	if err != nil {
		logger.Errorf("Error creating chat: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Could not create chat")
	}

	//一つ目のMessageを作成
	message := models.Message{
		ChatID:     chat.ID,
		UserUID:    userUID,
		Content:    "Hello",
		SenderType: "system",
	}

	_, err = h.messageService.CreateMessage(chat.ID, &message)
	if err != nil {
		logger.Errorf("Error creating message: %v", err)
		return respondWithError(c, http.StatusInternalServerError, "Could not create message")
	}

	logger.Infof("Chat created successfully")
	return c.JSON(http.StatusCreated, createdChat)

}

func (h *chatHandler) GetChatByMaterialID(c echo.Context) error {
	logger.Infof("Retrieving chat by material ID")
	materialID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid material ID")
	}

	userUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}

	chat, err := h.chatService.GetChatByMaterialID(materialID, userUID)
	if err != nil {
		logger.Errorf("Chat room not found: %v", err)
		return respondWithError(c, http.StatusNotFound, "Chat room not found")
	}

	logger.Infof("Chat retrieved successfully")
	return c.JSON(http.StatusOK, chat)
}

// GetChatByChatID retrieves a chat by its ID
func (h *chatHandler) GetChatByChatID(c echo.Context) error {
	id, err := parseUintParam(c, "chatId")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid chat ID")
	}

	userUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}

	chat, err := h.chatService.GetChatByChatID(id, userUID)
	if err != nil {
		logger.Errorf("Chat room not found: %v", err)
		return respondWithError(c, http.StatusNotFound, "Chat room not found")
	}

	logger.Infof("Chat retrieved successfully")
	return c.JSON(http.StatusOK, chat)
}

// ChatWithGemini handles communication with the Gemini API
func (h *chatHandler) ChatWithGemini(c echo.Context) error {
	chatID, err := parseUintParam(c, "chatId")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid chat ID")
	}

	userUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "Invalid user token")
	}

	var request struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&request); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request data")
	}

	response, err := h.messageService.SendMessageToGemini(chatID, request.Content, userUID)
	if err != nil {
		logger.Errorf("Error communicating with Gemini API: %v", err)
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"response": response})
}
