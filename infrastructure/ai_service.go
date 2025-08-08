package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiService wraps the Google Gemini generative AI client and model
type GeminiService struct {
	client *genai.Client          // Gemini API client instance
	model  *genai.GenerativeModel // Specific generative model to use
	ctx    context.Context         // Context for API calls
}

// NewGeminiService creates a new GeminiService given an API key
func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()

	// Initialize Gemini client with the provided API key
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		// Return error if client initialization fails
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	// Select the Gemini model to be used for generation
	model := client.GenerativeModel("models/gemini-1.5-flash")

	// Return the constructed GeminiService instance
	return &GeminiService{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}

// GenerateContent sends a prompt to the Gemini model and returns the generated text content
func (g *GeminiService) GenerateContent(prompt string) (string, error) {
	// Call the Gemini model's GenerateContent method with the prompt wrapped as genai.Text
	resp, err := g.model.GenerateContent(g.ctx, genai.Text(prompt))
	
	// Log the full response for debugging purposes
	log.Print("ZZZZZZZZZZZZZZZZZZZZZZZZ===", resp)

	// Return any errors encountered during the generation call
	if err != nil {
		return "", err
	}

	// Extract the generated text from the first candidate's first content part
	text := resp.Candidates[0].Content.Parts[0].(genai.Text)

	// Return the generated text as a string
	return string(text), nil
}
