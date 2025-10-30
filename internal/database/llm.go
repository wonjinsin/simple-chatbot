package database

import (
	"github.com/tmc/langchaingo/llms/ollama"
)

func NewOllamaLLM() (*ollama.LLM, error) {
	return ollama.New(
		ollama.WithModel("deepseek-r1:8b"),
	)
}
