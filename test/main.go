package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"seo-app/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	//	err := godotenv.Load("C:\\Users\\Никита\\Desktop\\КУРСЫ\\seo-app\\.env")
	////	if err != nil {
	////		log.Fatal("Error loading .env file")
	////	}
	////	apiKey := os.Getenv("AI_API_KEY")
	////	fmt.Println(apiKey)
	////	apiURL := os.Getenv("AI_API_URL")
	////	fmt.Println(apiURL)
	////	model := os.Getenv("AI_MODEL")
	////	fmt.Println(model)
	////
	////	const keywordExpansionPrompt = `Расширь ключевые слова для SEO.
	////Описание: "%s"
	////Базовые слова: %s
	////Сгенерируй 70-100 ключевых слов на русском.
	////Включи:
	////- Короткие (1-2 слова)
	////- Длинные (3-4 слова)
	////- С "купить", "цена"
	////- Геозависимые (Москва, СПб)
	////Формат: ["слово1", "слово2", ...]
	////
	////Только JSON массив.`
	////	aiService := service.NewAIService(apiKey, apiURL, model, nil)
	////	res, err := aiService.SendPrompt(keywordExpansionPrompt)
	////	if err != nil {
	////		fmt.Println(err)
	////	}
	////	fmt.Println(res)

	if err := godotenv.Load("C:\\Users\\Никита\\Desktop\\КУРСЫ\\seo-app\\.env"); err != nil {
		log.Println("No .env file found, using system environment")
	}
	wordstatApiKey := os.Getenv("WORDSTAT_API_KEY")
	baseURL := os.Getenv("WORDSTAT_API_URL")
	wsService := service.NewWordstatService(baseURL, wordstatApiKey, nil)
	req := service.TopRequestsRequest{
		Phrases: []string{"магазин часов", "купить часы мужские", "часы casio"},
	}

	result, err := wsService.RequestWordstat(&req)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
