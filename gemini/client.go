package gemini

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var geminiClient *genai.Client

func init() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("WARNING: GEMINI_API_KEY environment variable not set")
		return
	}

	initClient(apiKey)
}

func initClient(apiKey string) error {
	var err error
	geminiClient, err = genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %v", err)
	}
	return nil
}

func GetGeminiResponse(ctx context.Context, prompt string) (string, error) {
	if geminiClient == nil {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("Gemini client not initialized, please set GEMINI_API_KEY environment variable")
		}
		
		if err := initClient(apiKey); err != nil {
			return "", err
		}
	}

	model := geminiClient.GenerativeModel("gemini-2.0-flash")
	
	model.SetTemperature(0.2)
	model.SetTopP(0.8)
	model.SetTopK(40)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}