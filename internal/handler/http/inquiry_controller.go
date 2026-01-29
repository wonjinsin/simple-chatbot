package http

import (
	"net/http"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
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
		logger.LogError(ctx, "EmbedInquiryOrigins failed", err)
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

// Ask handles inquiry request and returns the most similar domain
func (c *InquiryController) Ask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "Ask request received")

	// Step 1: Parse request body
	var req dto.AskRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		logger.LogWarn(ctx, "invalid json in request body")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "invalid json",
		}, string(constants.InvalidParameter))
		return
	}

	// Step 2: Call service to find the most similar inquiry knowledge
	result, err := c.svc.Ask(ctx, req.Msg)
	if err != nil {
		logger.LogError(ctx, "Ask failed", err)
		// Extract error code and determine HTTP status
		code := errors.GetCode(err)

		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(code))
		return
	}

	// Step 3: Convert domain to DTO
	response := dto.ToAskResponse(result)

	logger.LogInfo(ctx, "Ask success response received")
	utils.WriteStandardJSON(w, r, http.StatusOK, response)
}
