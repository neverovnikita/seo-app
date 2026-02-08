package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"seo-app/internal/dto"
	"time"
)

type AIService struct {
	apiKey     string
	apiURL     string
	model      string
	httpClient *http.Client
}

type AIRequest struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type KeywordsResult struct {
	Keywords []string `json:"keywords"`
	Count    int      `json:"count"`
}

func NewAIService(apiKey string, apiURL string, model string, httpClient *http.Client) *AIService {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	}
	return &AIService{apiKey, apiURL, model, httpClient}
}

//func (s *AIService) SendMessage(input *dto.ProjectResponse) (string, error) {
//	request := AIRequest{
//		Model: s.model,
//		Messages: []AIMessage{
//			{
//				Role:    "user",
//				Content: s.prompt(input),
//			},
//		},
//	}
//	jsonData, err := json.Marshal(request)
//	if err != nil {
//		return "", fmt.Errorf("failed to marshal request: %w", err)
//	}
//	req, err := http.NewRequest(
//		"POST",
//		"https://openrouter.ai/api/v1/chat/completions",
//		bytes.NewBuffer(jsonData),
//	)
//	if err != nil {
//		return "", fmt.Errorf("failed to create request: %w", err)
//	}
//	req.Header.Set("Authorization", "Bearer "+s.apiKey)
//	req.Header.Set("Content-Type", "application/json")
//
//	client := s.httpClient
//	resp, err := client.Do(req)
//	log.Println("отправка запроса")
//
//	if err != nil {
//		return "", fmt.Errorf("HTTP request failed: %w", err)
//	}
//	defer resp.Body.Close()
//	body, err := io.ReadAll(resp.Body)
//
//	if err != nil {
//		return "", fmt.Errorf("failed to read response: %w", err)
//	}
//	if resp.StatusCode != http.StatusOK {
//		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
//	}
//
//	var aiResp AIResponse
//	if err := json.Unmarshal(body, &aiResp); err != nil {
//		return "", fmt.Errorf("failed to parse response: %w", err)
//	}
//
//	if len(aiResp.Choices) == 0 {
//		return "", fmt.Errorf("no response from AI")
//	}
//
//	log.Println("конец метода")
//	log.Println(aiResp.Choices[0].Message.Content)
//
//	return aiResp.Choices[0].Message.Content, nil
//}

func (s *AIService) prompt(input *dto.ProjectResponse) string {
	const keywordExpansionPrompt = `Расширь ключевые слова для SEO.
Описание: "%s"
Базовые слова: %s
Сгенерируй 70-100 ключевых слов на русском.
Включи: 
- Короткие (1-2 слова)
- Длинные (3-4 слова)  
- С "купить", "цена"
- Геозависимые (Москва, СПб)
Формат: ["слово1", "слово2", ...]
`
	keywordsJSON, _ := json.Marshal(input.BaseKeywords)
	return fmt.Sprintf(keywordExpansionPrompt, input.Description, string(keywordsJSON))
}
func (s *AIService) SendPrompt(input *dto.ProjectResponse) (KeywordsResult, error) {
	prompt := s.prompt(input)
	log.Println(prompt)
	request := AIRequest{
		Model: s.model,
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	result := KeywordsResult{}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return result, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := s.httpClient
	resp, err := client.Do(req)
	log.Println("отправка запроса")
	if err != nil {
		return result, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return result, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return result, fmt.Errorf("no response from AI")
	}

	log.Println("конец метода")
	log.Println(aiResp.Choices[0].Message.Content)

	var keywords []string
	if err := json.Unmarshal([]byte(aiResp.Choices[0].Message.Content), &keywords); err == nil {
		result.Keywords = keywords
		result.Count = len(keywords)
	}

	return result, nil
}
