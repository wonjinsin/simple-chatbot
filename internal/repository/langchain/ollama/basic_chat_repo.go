package langchain

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
)

type basicChatRepo struct {
	ollamaLLM *ollama.LLM
}

// NewBasicChatRepo creates a new basic chat repository
func NewBasicChatRepo(ollamaLLM *ollama.LLM) repository.BasicChatRepository {
	return &basicChatRepo{ollamaLLM: ollamaLLM}
}

// Ask asks the LLM a question and returns the answer
func (r *basicChatRepo) Ask(ctx context.Context, msg string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, r.ollamaLLM, msg)
	if err != nil {
		return "", err
	}

	return completion, nil
}
