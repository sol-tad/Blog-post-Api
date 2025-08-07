package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)
const geminiURL = "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent"

type GeminiRequest struct{
	Contents []Content `json:"contents"`
}

type Content struct{
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct{
	Candidates []struct{
		Content struct{
			Parts[]struct{
				Text string `json:"text`
			} `json:"parts"`
		} `json:"content"`
	}`json:"candidates"`
}

func SummerizeBlog(content string) (string, error){
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("%s?key=%s", geminiURL, apiKey)

	prompt := fmt.Sprintf("Summarize the following blog post:\n\n%s", content)	
	reqBody := GeminiRequest{
		Contents: []Content{
			{Parts: []Part{{Text:prompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)



	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}
	if len(geminiResp.Candidates) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response from Gemini")

}

func EmphasizeBlog(content string) (string, error){
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("%s?key=%s", geminiURL, apiKey)

	prompt := fmt.Sprintf("emphasize more on the following blog post:\n\n%s", content)	
	reqBody := GeminiRequest{
		Contents: []Content{
			{Parts: []Part{{Text:prompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)



	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}
	if len(geminiResp.Candidates) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response from Gemini")

}


func TitleBlog(content string) (string, error){
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("%s?key=%s", geminiURL, apiKey)

	prompt := fmt.Sprintf("give me best title for the following blog post:\n\n%s", content)	
	reqBody := GeminiRequest{
		Contents: []Content{
			{Parts: []Part{{Text:prompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)



	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}
	if len(geminiResp.Candidates) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response from Gemini")

}




func ImproveBlog(content string) (string, error){
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("%s?key=%s", geminiURL, apiKey)

	prompt := fmt.Sprintf("improve the following blog post:\n\n%s", content)	
	reqBody := GeminiRequest{
		Contents: []Content{
			{Parts: []Part{{Text:prompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)



	var geminiResp GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}
	if len(geminiResp.Candidates) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no response from Gemini")

}

