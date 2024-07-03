package gemini

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/yomek33/talki/internal/models"
)

// SendMessageToGemini sends a message to the Gemini model and returns the response
func (c *Client) SendMessageToGemini(ctx context.Context, chat *models.Chat, content string) (string, error) {
	geminiModel := c.client.GenerativeModel("gemini-1.5-flash")
	cs := geminiModel.StartChat()

	// Convert chat messages to Gemini API format
	for _, msg := range chat.Messages {
		role := "user"
		if msg.SenderType == "bot" {
			role = "model"
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(msg.Content),
			},
			Role: role,
		})
	}

	// Append the current user message
	cs.History = append(cs.History, &genai.Content{
		Parts: []genai.Part{
			genai.Text(content),
		},
		Role: "user",
	})

	// Send the message to Gemini
	resp, err := cs.SendMessage(ctx, genai.Text(content))
	if err != nil {
		log.Printf("Error sending message to Gemini: %v", err)
		return "", fmt.Errorf("error sending message to Gemini: %w", err)
	}

	// Validate response
	if resp.Candidates == nil || len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	if resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	// Type assertion with check
	part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected content type in response")
	}

	// Print response for debugging
	printResponse(resp)

	return string(part), nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Printf("Part: %v\n", part)
			}
		}
	}
}
