package gemini

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

var GeminiClient *genai.Client

func InitGemini() error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	GeminiClient = client
	return nil
}
