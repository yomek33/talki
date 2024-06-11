package gemini_test

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/joho/godotenv"
// 	"github.com/yomek33/talki/internal/gemini"
// )

// func TestGenerateJsonContent(t *testing.T) {
// 	ctx := context.Background() // Replace with a valid context

// 	if err := godotenv.Load(); err != nil {
// 		t.Fatal(err)
// 	}
// 	apiKey := os.Getenv("GEMINI_API_KEY")

// 	fmt.Println(apiKey)
// 	client, err := gemini.NewClient(ctx, apiKey)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer client.Close()

// 	t.Run("successful generation", func(t *testing.T) {
// 		prompt := "What is the weather like today?"
// 		output, err := client.GenerateJsonContent(ctx, prompt)
// 		if err != nil {
// 			t.Errorf("unexpected error: %v", err)
// 		}
// 		if output == nil || len(output) == 0 {
// 			t.Errorf("expected non-empty output")
// 		}
// 		for _, o := range output {
// 			fmt.Println(o)
// 		}
// 	})
// }
