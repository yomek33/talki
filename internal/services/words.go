package services

import "github.com/yomek33/talki/internal/gemini"

type WordService interface {
	GenerateWords(topic string) ([]string, error)
}


type wordService struct {
	GeminiClient *gemini.Client
}


