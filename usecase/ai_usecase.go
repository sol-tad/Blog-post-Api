package usecase
import (
	"github.com/sol-tad/Blog-post-Api/infrastructure"
)
type AiUseCase struct { }

func NewAiUseCase() AiUseCase{
	return AiUseCase{}
}

func (au *AiUseCase) GenerateSummary (content string) (string , error){
	return infrastructure.SummerizeBlog(content)
}

func (au *AiUseCase) GenerateImproved(content string) (string, error){
	return infrastructure.ImproveBlog(content)
}

func (au *AiUseCase) GenerateEmphasized(content string) (string, error){
	return infrastructure.EmphasizeBlog(content)
}

func (au *AiUseCase) GenerateTitle(content string) (string, error){
	return infrastructure.TitleBlog(content)
}