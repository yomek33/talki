package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/yomek33/talki/internal/models"
)

func (c *Client) SendMessageToGemini(ctx context.Context, chat *models.Chat, content string) ([]string, error) {
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

	cs.History = append(cs.History, &genai.Content{
		Parts: []genai.Part{
			genai.Text(content),
		},
		Role: "user",
	})

	resp, err := cs.SendMessage(ctx, genai.Text(content))
	if err != nil {
		return nil, err
	}
	return printResponse(resp), nil
}

func printResponse(resp *genai.GenerateContentResponse) []string {
	var meassages []string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				meassages = append(meassages, fmt.Sprintf("%v", part))
			}
		}
	}
	return meassages
}
