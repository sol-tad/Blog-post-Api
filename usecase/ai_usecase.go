package usecase

import (
	"fmt"
	"log"

	"github.com/sol-tad/Blog-post-Api/domain"
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)

type AIUsecase struct {
	gemini *infrastructure.GeminiService
}
func NewAIUsecase(gemini *infrastructure.GeminiService) domain.AIUsecase{
	return &AIUsecase{
		gemini:gemini,
	}
}

//genretae blog post
func(aiu *AIUsecase) GenerateBlogPost(input domain.GenerateBlogPostRequest)(string,error){
	prompt:=fmt.Sprintf("Write a %s-lenght blog post about '%s' for %s in a %s tone.",input.Length,input.Topic,input.Audience,input.Tone)
	log.Print("444444444444444444444444444",prompt)

return aiu.gemini.GenerateContent(prompt)
}

//imporovw blog post
func(aiu *AIUsecase) ImproveBlogPost(input domain.ImproveBlogPostRequest) (string, error) {
	prompt := fmt.Sprintf("Improve the following blog post with the goal: %s\n\n%s", input.Goal, input.Content)
	return aiu.gemini.GenerateContent(prompt)

}


//suggest Improvemnets
func (aiu *AIUsecase) SuggestBlogImprovement(content string) ([]string, error) {

	prompt:=fmt.Sprintf("Suggest writing and SEO improvements for this blog post:\n\n%s", content)
		resp, err := aiu.gemini.GenerateContent(prompt)

		if err != nil {
		return nil, err
	}
	return []string{resp}, nil
}

//Summarize

func (aiu *AIUsecase) SummarizeBlog(content string) (string, error) {
	prompt := fmt.Sprintf("Summarize the following blog post in 2-3 sentences:\n\n%s", content)
	return aiu.gemini.GenerateContent(prompt)
}

func (aiu *AIUsecase) GenerateMetadata(content string) (map[string]string, error) {
	prompt := fmt.Sprintf("Based on the following blog post, suggest a better title, 5 tags, and a meta description:\n\n%s", content)
	result, err := aiu.gemini.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"metadata": result,
	}, nil
}