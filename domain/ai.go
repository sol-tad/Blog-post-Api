package domain

// GenerateBlogPostRequest represents the input data for generating a new blog post using AI.
// swagger:model GenerateBlogPostRequest
type GenerateBlogPostRequest struct {
	// The main topic or subject of the blog post
	// example: "Introduction to Go concurrency"
	Topic string `json:"topic"`

	// The desired tone or style of the writing (e.g., formal, casual)
	// example: "formal"
	Tone string `json:"tone"`

	// The preferred length of the blog post (e.g., short, medium, long)
	// example: "medium"
	Length string `json:"length"`

	// The target audience for the blog post (e.g., beginners, experts)
	// example: "beginners"
	Audience string `json:"audience"`
}

// ImproveBlogPostRequest represents the input data for improving an existing blog post.
// swagger:model ImproveBlogPostRequest
type ImproveBlogPostRequest struct {
	// The current content of the blog post to be improved
	// example: "Go is a statically typed programming language..."
	Content string `json:"content"`

	// The goal or focus of the improvement (e.g., clarity, engagement)
	// example: "increase clarity"
	Goal string `json:"goal"`
}

// SuggestionResponse represents the response structure containing AI suggestions for improvements.
// swagger:model SuggestionResponse
type SuggestionResponse struct {
	// List of suggestions
	// example: ["Use simpler language", "Add examples"]
	Suggestions []string `json:"suggestions"`  // fixed typo from "suggesions"
}

// AIUsecase defines the interface for AI-related blog operations.
// No swagger annotation needed for interfaces
type AIUsecase interface {
	GenerateBlogPost(input GenerateBlogPostRequest) (string, error)           // Generates a new blog post based on input parameters
	ImproveBlogPost(input ImproveBlogPostRequest) (string, error)            // Improves an existing blog post content based on a goal
	SuggestBlogImprovement(content string) ([]string, error)                 // Provides suggestions to improve the blog content
	SummarizeBlog(content string) (string, error)                            // Generates a summary of the blog content
	GenerateMetadata(content string) (map[string]string, error)              // Generates metadata (e.g., tags, keywords) for the blog content
}

