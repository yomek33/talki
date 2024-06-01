package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// GenerateContent generates content based on the provided prompt
func (c *Client) GenerateJsonContent(ctx context.Context, prompt string) ([]string, error) {
	model := c.client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}

	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	var output []string
	for _, part := range res.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			outputText := strings.ReplaceAll(string(txt), "\\", "")
			outputText = strings.ReplaceAll(outputText, "\"", "")
			output = append(output, outputText)
		}
	}

	return output, nil
}

// GeneratePhrases generates phrases based on the given topic
func (c *Client)GeneratePhrases(ctx context.Context, topic string) ([]string, error) {
	prompt := generatePrompt(topic)
	output, err := c.GenerateJsonContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return output, nil
}

func generatePrompt(topic string) string {
	promptParts := []string{
		"Generate 10 useful English phrases related to {topic}, focusing on {action verb} (e.g., describing, discussing). Include synonyms and related terms for {topic}.",
		"topic: climate change",
		"output: [ \"The planet is experiencing an unprecedented rise in global temperatures.\", \"Human activities are the primary drivers of climate change.\", \"Rising sea levels threaten coastal communities around the world.\", \"Extreme weather events, such as hurricanes and heatwaves, are becoming more frequent and intense.\", \"Greenhouse gases, such as carbon dioxide and methane, trap heat in the atmosphere.\", \"Climate change poses a significant threat to biodiversity and ecosystems.\", \"Renewable energy sources, such as solar and wind power, are essential for mitigating climate change.\", \"Carbon emissions must be drastically reduced to limit global warming.\", \"Climate change is a complex and urgent issue that requires global cooperation.\", \"Sustainable practices, such as reducing consumption and improving energy efficiency, are crucial for addressing climate change.\" ]",
		fmt.Sprintf("topic: %s", topic),
		"output: ",
	}

	return strings.Join(promptParts, "\n")
}
