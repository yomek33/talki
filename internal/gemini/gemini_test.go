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
			{Content: "こんにちは。私は広島に住んでいるので聞いてください", SenderType: "user"},
			{Content: "こんにちは!あなたは広島に住んでいるのですね", SenderType: "bot"},
		},
	}

	response, err := client.SendMessageToGemini(ctx, chat, "私に何を聞きたいですか")
	fmt.Println(response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)

}
