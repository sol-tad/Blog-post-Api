package domain

type GenerateBlogPostRequest struct {
	Topic    string `json:"topic"`
	Tone     string `json:"tone"`
	Length   string `json:"lenght"`
	Audience string `json:"audience"`
}

type ImproveBlogPostRequest struct {
	Content string `json:"content"`
	Goal    string `json:"goal"`
}

type SuggestionResponse struct {
	Suggestion []string `json:"suggesions"`
}

type AIUsecase interface {
	GenerateBlogPost(input GenerateBlogPostRequest) (string, error)
	ImproveBlogPost(input ImproveBlogPostRequest) (string, error)
	SuggestBlogImprovement(content string) ([]string, error)
	SummarizeBlog(content string)(string,error)
}