package infrastructure

import (
	"context"
	"fmt"
	"log"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	model := client.GenerativeModel("models/gemini-1.5-flash")

	return &GeminiService{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}


func (g *GeminiService) GenerateContent(prompt string) (string, error) {
	resp, err := g.model.GenerateContent(g.ctx, genai.Text(prompt))
	log.Print("ZZZZZZZZZZZZZZZZZZZZZZZZ===",resp)
	if err != nil {
		return "", err
	}

text := resp.Candidates[0].Content.Parts[0].(genai.Text)
return string(text), nil
}
