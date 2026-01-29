package dto

import "github.com/wonjinsin/simple-chatbot/internal/domain"

// ToAskResponse converts InquirySimilarityResult domain object to AskResponse DTO
func ToAskResponse(result *domain.InquirySimilarityResult) *AskResponse {
	if result == nil || result.Knowledge == nil {
		return nil
	}

	return &AskResponse{
		Response: result.Knowledge.Response,
	}
}
