package usecase

import "github.com/sol-tad/Blog-post-Api/domain"


type AIUseCase struct {
	aiService domain.AIService
}

func NewAIUseCases(aiService domain.AIService) *AIUseCase {
	return &AIUseCase{
		aiService: aiService,
	}}

func (ac *AIUseCase) GenerateBlog(params domain.GenerationParams) ( string, error){
	if params.Tone == ""{
		params.Tone = "professional"
	}
	if params.Length == 0{
		params.Length = 800
	}
	return ac.aiService.GenerateContent(params)
}

func(ac *AIUseCase) SummarizeBLog(content string) (string, error){
	return ac.SummarizeBLog(content)
}