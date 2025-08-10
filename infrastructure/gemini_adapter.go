package infrastructure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sol-tad/Blog-post-Api/domain"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiAdapter struct {
	model *genai.GenerativeModel
}

func NewGeminiAdapter() (domain.AIService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Gemini client error: %w", err)
	}

	return &GeminiAdapter{
		model: client.GenerativeModel("models/gemini-1.5-flash"),
	}, nil
}

func (g *GeminiAdapter) GenerateContent(params domain.GenerationParams) (string, error) {
	ctx := context.Background()

	// Build detailed prompt
	prompt := buildPrompt(params)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return extractResponse(resp), nil
}

func (g *GeminiAdapter) SummarizeBlog(content string) (string, error) {
	prompt := fmt.Sprintf(
		"Summarize this blog post in 3 bullet points:\n\n%s",
		content,
	)
	resp, err := g.model.GenerateContent(context.Background(), genai.Text(prompt))
	return extractResponse(resp), err
}

// Helper to build sophisticated prompt
func buildPrompt(params domain.GenerationParams) string {
	var sb strings.Builder

	sb.WriteString("Write a comprehensive blog post with these requirements:\n")
	sb.WriteString(fmt.Sprintf("- Topic: %s\n", params.Topic))

	if params.Title != "" {
		sb.WriteString(fmt.Sprintf("- Title: %s\n", params.Title))
	}

	if len(params.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("- Tags: %s\n", strings.Join(params.Tags, ", ")))
	}

	sb.WriteString(fmt.Sprintf("- Tone: %s\n", params.Tone))
	sb.WriteString(fmt.Sprintf("- Length: %d words\n", params.Length))

	sb.WriteString("\nStructure:\n")
	sb.WriteString("1. Engaging introduction\n")
	sb.WriteString("2. 3-5 main sections with subheadings\n")
	sb.WriteString("3. Conclusion with key takeaways\n")

	return sb.String()
}

func extractResponse(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	}
	return ""
}
