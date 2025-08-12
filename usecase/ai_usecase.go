package usecase

import (
	"fmt"
	"log"

	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

// AIUsecase handles AI-powered blog operations using GeminiService
type AIUsecase struct {
	gemini *infrastructure.GeminiService
}

// NewAIUsecase initializes AIUsecase with a GeminiService instance
func NewAIUsecase(gemini *infrastructure.GeminiService) domain.AIUsecase {
	return &AIUsecase{
		gemini: gemini,
	}
}

// GenerateBlogPost creates a new blog post based on user input
func (aiu *AIUsecase) GenerateBlogPost(input domain.GenerateBlogPostRequest) (string, error) {
	// Build prompt using input parameters: length, topic, audience, tone
	prompt := fmt.Sprintf(
		"Write a %s-length blog post about '%s' for %s in a %s tone.",
		input.Length, input.Topic, input.Audience, input.Tone,
	)

	log.Print("Generated prompt for blog post:", prompt)

	// Send prompt to GeminiService to generate content
	return aiu.gemini.GenerateContent(prompt)
}

// ImproveBlogPost enhances an existing blog post based on a specific goal
func (aiu *AIUsecase) ImproveBlogPost(input domain.ImproveBlogPostRequest) (string, error) {
	// Build prompt with improvement goal and original content
	prompt := fmt.Sprintf(
		"Improve the following blog post with the goal: %s\n\n%s",
		input.Goal, input.Content,
	)

	// Generate improved content via GeminiService
	return aiu.gemini.GenerateContent(prompt)
}

// SuggestBlogImprovement provides writing and SEO suggestions for a blog post
func (aiu *AIUsecase) SuggestBlogImprovement(content string) ([]string, error) {
	// Build prompt to request suggestions
	prompt := fmt.Sprintf(
		"Suggest writing and SEO improvements for this blog post:\n\n%s",
		content,
	)

	// Get suggestions from GeminiService
	resp, err := aiu.gemini.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	// Return suggestions as a slice (can be expanded later)
	return []string{resp}, nil
}

// SummarizeBlog generates a short summary (2–3 sentences) of the blog post
func (aiu *AIUsecase) SummarizeBlog(content string) (string, error) {
	prompt := fmt.Sprintf(
		"Summarize the following blog post in 2-3 sentences:\n\n%s",
		content,
	)

	return aiu.gemini.GenerateContent(prompt)
}

// GenerateMetadata extracts SEO metadata like title, tags, and description
func (aiu *AIUsecase) GenerateMetadata(content string) (map[string]string, error) {
	prompt := fmt.Sprintf(
		"Based on the following blog post, suggest a better title, 5 tags, and a meta description:\n\n%s",
		content,
	)

	result, err := aiu.gemini.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	// Return metadata as a simple map (can be parsed into structured fields later)
	return map[string]string{
		"metadata": result,
	}, nil
}