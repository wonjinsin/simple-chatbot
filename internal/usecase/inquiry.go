package usecase

import (
	"context"
	"fmt"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/file"
)

type inquiryService struct {
}

func NewInquiryService() InquiryService {
	return &inquiryService{}
}

// EmbedInquiryOrigins embeds inquiry origins
func (s *inquiryService) EmbedInquiryOrigins(_ context.Context) error {
	origins, err := file.ReadCSVToMapArray("mock_data/data_set.csv")
	if err != nil {
		return errors.Wrap(err, "failed to read inquiry origins", constants.InternalError)
	}

	for _, origin := range origins {
		fmt.Println(origin)
		break
	}
	return nil
}
