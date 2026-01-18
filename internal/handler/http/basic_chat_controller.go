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

// BasicChatController handles basic chat-related HTTP requests
type BasicChatController struct {
	svc usecase.BasicChatService
}

// NewBasicChatController creates a new basic chat controller
func NewBasicChatController(svc usecase.BasicChatService) *BasicChatController {
	return &BasicChatController{svc: svc}
}

// Ask handles basic chat request
func (c *BasicChatController) AskBasicChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "ask request received")

	var req dto.AskRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		logger.LogWarn(ctx, "invalid json in request body")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "invalid json",
		}, string(constants.InvalidParameter))
		return
	}

	answer, err := c.svc.AskBasicChat(ctx, req.Msg)
	if err != nil {
		logger.LogError(ctx, "internal error in ask", err)
		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(errors.GetCode(err)))
		return
	}

	logger.LogInfo(ctx, "answer retrieved successfully")
	utils.WriteStandardJSON(w, r, http.StatusOK, answer)
}

// AskBasicPromptTemplateChat handles basic prompt template chat request
func (c *BasicChatController) AskBasicPromptTemplateChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "ask basic prompt template chat request received")

	var req dto.AskRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		logger.LogWarn(ctx, "invalid json in request body")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "invalid json",
		}, string(constants.InvalidParameter))
		return
	}

	answer, err := c.svc.AskBasicPromptTemplateChat(ctx, req.Msg)
	if err != nil {
		logger.LogError(ctx, "internal error in ask basic prompt template chat", err)
		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(errors.GetCode(err)))
		return
	}

	logger.LogInfo(ctx, "answer basic prompt template chat retrieved successfully")
	utils.WriteStandardJSON(w, r, http.StatusOK, answer)
}
