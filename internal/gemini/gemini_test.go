package gemini

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// func TestSendMessageToGemini(t *testing.T) {
// 	ctx := context.Background()
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Printf("error loading .env file: %v", err)
// 	}
// 	// Initialize the real API client
// 	apiKey := os.Getenv("GEMINI_API_KEY")
// 	if apiKey == "" {
// 		t.Skip("GENAI_API_KEY environment variable is not set")
// 	}
// 	client, err := NewClient(ctx, apiKey)
// 	if err != nil {
// 		t.Fatalf("failed to create client: %v", err)
// 	}
// 	defer client.Close()
// 	promptParts := []string{
// 		"Generate 10 useful English phrases related to {topic}, focusing on {action verb} (e.g., describing, discussing). Include synonyms and related terms for {topic}.",
// 		"topic: climate change",
// 		"output: [ \"The planet is experiencing an unprecedented rise in global temperatures.\", \"Human activities are the primary drivers of climate change.\", \"Rising sea levels threaten coastal communities around the world.\", \"Extreme weather events, such as hurricanes and heatwaves, are becoming more frequent and intense.\", \"Greenhouse gases, such as carbon dioxide and methane, trap heat in the atmosphere.\", \"Climate change poses a significant threat to biodiversity and ecosystems.\", \"Renewable energy sources, such as solar and wind power, are essential for mitigating climate change.\", \"Carbon emissions must be drastically reduced to limit global warming.\", \"Climate change is a complex and urgent issue that requires global cooperation.\", \"Sustainable practices, such as reducing consumption and improving energy efficiency, are crucial for addressing climate change.\" ]",
// 		fmt.Sprintf("topic: %s", topic),
// 		"output: ",
// 	}
// 	prom = string(promptParts)
// 	chat := &models.Chat{
// 		Detail: "Example Chat",
// 		Messages: []models.Message{
// 			{Content: promptParts, SenderType: "user"},
// 			{Content: "こんにちは!あなたは広島に住んでいるのですね", SenderType: "bot"},
// 		},
// 	}

// 	response, err := client.SendMessageToGemini(ctx, chat, "私に何を聞きたいですか")
// 	fmt.Println(response)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, response)

// }


func TestGenerateWords(t *testing.T) {
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

	inWords, err := client.GenerateIntermediateWords(ctx, "A large language model (LLM) is a computational model notable for its ability to achieve general-purpose language generation and other natural language processing tasks such as classification. Based on language models, LLMs acquire these abilities by learning statistical relationships from vast amounts of text during a computationally intensive self-supervised and semi-supervised training process.[1] LLMs can be used for text generation, a form of generative AI, by taking an input text and repeatedly predicting the next token or word.[2]")
	fmt.Println(inWords)
	assert.NoError(t, err)
	assert.NotEmpty(t, inWords)

	adWords, err := client.GenerateAdvancedWords(ctx, "A large language model (LLM) is a computational model notable for its ability to achieve general-purpose language generation and other natural language processing tasks such as classification. Based on language models, LLMs acquire these abilities by learning statistical relationships from vast amounts of text during a computationally intensive self-supervised and semi-supervised training process.[1] LLMs can be used for text generation, a form of generative AI, by taking an input text and repeatedly predicting the next token or word.[2]")
	fmt.Println(adWords)
	assert.NoError(t, err)
	assert.NotEmpty(t, adWords)
}