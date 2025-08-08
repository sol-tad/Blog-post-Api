package domain

// GenerateBlogPostRequest represents the input data for generating a new blog post using AI.
type GenerateBlogPostRequest struct {
	Topic    string `json:"topic"`    // The main topic or subject of the blog post
	Tone     string `json:"tone"`     // The desired tone or style of the writing (e.g., formal, casual)
	Length   string `json:"length"`   // The preferred length of the blog post (e.g., short, medium, long)
	Audience string `json:"audience"` // The target audience for the blog post (e.g., beginners, experts)
}

// ImproveBlogPostRequest represents the input data for improving an existing blog post.
type ImproveBlogPostRequest struct {
	Content string `json:"content"` // The current content of the blog post to be improved
	Goal    string `json:"goal"`    // The goal or focus of the improvement (e.g., clarity, engagement)
}

// SuggestionResponse represents the response structure containing AI suggestions for improvements.
type SuggestionResponse struct {
	Suggestion []string `json:"suggesions"` // List of suggestions (note: typo "suggesions" should be "suggestions")
}

// AIUsecase defines the interface for AI-related blog operations.
type AIUsecase interface {
	GenerateBlogPost(input GenerateBlogPostRequest) (string, error)           // Generates a new blog post based on input parameters
	ImproveBlogPost(input ImproveBlogPostRequest) (string, error)            // Improves an existing blog post content based on a goal
	SuggestBlogImprovement(content string) ([]string, error)                 // Provides suggestions to improve the blog content
	SummarizeBlog(content string) (string, error)                            // Generates a summary of the blog content
	GenerateMetadata(content string) (map[string]string, error)              // Generates metadata (e.g., tags, keywords) for the blog content
}
