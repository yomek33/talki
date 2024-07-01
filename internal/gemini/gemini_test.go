package gemini

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/yomek33/talki/internal/models"
)

func TestSendMessageToGemini(t *testing.T) {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}
	// Initialize the real API client
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GENAI_API_KEY environment variable is not set")
	}
	client, err := NewClient(ctx, apiKey)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	chat := &models.Chat{
		Detail: "Example Chat",
		Messages: []models.Message{
			{Content: "Hello, my job is a software engineer", SenderType: "user"},
			{Content: "you are a software engineer!", SenderType: "bot"},
		},
	}

	response, err := client.SendMessageToGemini(ctx, chat, "what is my job?")
	fmt.Println(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)

	for _, msg := range response {
		t.Log(msg)
	}
}
