package usecase

import (
	"context"

	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type basicChatService struct {
	repo repository.BasicChatRepository
}

func NewBasicChatService(r repository.BasicChatRepository) BasicChatService {
	return &basicChatService{repo: r}
}

func (s *basicChatService) AskBasicChat(ctx context.Context, msg string) (string, error) {
	answer, err := s.repo.AskBasicChat(ctx, msg)
	if err != nil {
		return "", errors.Wrap(err, "failed to ask")
	}
	return answer, nil
}

func (s *basicChatService) AskBasicPromptTemplateChat(ctx context.Context, msg string) (string, error) {
	answer, err := s.repo.AskWithTool(ctx, msg)
	if err != nil {
		return "", errors.Wrap(err, "failed to ask")
	}
	return answer, nil
}
