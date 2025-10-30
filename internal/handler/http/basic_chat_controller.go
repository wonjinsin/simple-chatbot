package http

import (
	"fmt"
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
func (c *BasicChatController) Ask(w http.ResponseWriter, r *http.Request) {
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

	answer, err := c.svc.Ask(ctx, req.Msg)
	if err != nil {
		logger.LogError(ctx, "internal error in ask", err)

		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(errors.GetCode(err)))
		return
	}

	logger.LogInfo(ctx, fmt.Sprintf("answer retrieved successfully: %s", answer))

	utils.WriteStandardJSON(w, r, http.StatusOK, answer)
}
