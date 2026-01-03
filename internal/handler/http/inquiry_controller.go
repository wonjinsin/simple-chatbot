package http

import (
	"net/http"

	"github.com/wonjinsin/simple-chatbot/internal/handler/http/dto"
	"github.com/wonjinsin/simple-chatbot/internal/usecase"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/logger"
	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

// InquiryController handles inquiry-related HTTP requests
type InquiryController struct {
	svc usecase.InquiryService
}

// NewInquiryController creates a new inquiry controller
func NewInquiryController(svc usecase.InquiryService) *InquiryController {
	return &InquiryController{svc: svc}
}

// EmbedInquiryOrigins handles inquiry origins embedding
func (c *InquiryController) EmbedInquiryOrigins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "EmbedInquiryOrigins request received")

	err := c.svc.EmbedInquiryOrigins(ctx)
	if err != nil {
		// Extract error code and determine HTTP status
		code := errors.GetCode(err)

		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(code))
		return
	}

	logger.LogInfo(ctx, "EmbedInquiryOrigins success response received")
	utils.WriteStandardJSON(w, r, http.StatusCreated, "success")
}
