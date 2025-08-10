package domain

// AI Generation Parameters
type GenerationParams struct {
    Topic     string   `json:"topic"`
    Title     string   `json:"title,omitempty"`
    Tags      []string `json:"tags,omitempty"`
    Tone      string   `json:"tone,omitempty"` // e.g., "professional", "casual", "humorous"
    Length    int      `json:"length,omitempty"` // word count
}

// AI Response Structure
type AIResponse struct {
    Content string `json:"content"`
}

// AI Service Interface
type AIService interface {
    GenerateContent(params GenerationParams) (string, error)
    SummarizeBlog(content string) (string, error)
}