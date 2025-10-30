package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/handler/http/dto"
	"github.com/wonjinsin/simple-chatbot/internal/usecase"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
	"github.com/wonjinsin/simple-chatbot/pkg/logger"
	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

// UserController handles user-related HTTP requests
type UserController struct {
	svc usecase.UserService
}

// NewUserController creates a new user controller
func NewUserController(svc usecase.UserService) *UserController {
	return &UserController{svc: svc}
}

// CreateUser handles user creation
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "CreateUser request received")

	var req dto.CreateUserRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		logger.LogWarn(ctx, "invalid json in request body")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "invalid json",
		}, string(constants.InvalidParameter))
		return
	}

	u, err := c.svc.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		// Extract error code and determine HTTP status
		code := errors.GetCode(err)
		var httpStatus int

		switch code {
		case constants.ConstraintError:
			httpStatus = http.StatusConflict
			logger.LogWarn(ctx, "duplicate email attempted")
		case constants.InvalidParameter:
			httpStatus = http.StatusBadRequest
			logger.LogWarn(ctx, "invalid parameter in create user")
		case "":
			// Non-CustomError, default to 500
			code = constants.InternalError
			httpStatus = http.StatusInternalServerError
			logger.LogError(ctx, "unexpected error in create user", err)
		default:
			httpStatus = http.StatusInternalServerError
			logger.LogError(ctx, "internal error in create user", err)
		}

		utils.WriteStandardJSON(w, r, httpStatus, dto.ErrorResult{
			Msg: err.Error(),
		}, string(code))
		return
	}

	logger.LogInfo(ctx, "user created successfully")
	response := dto.ToUserResponse(u)
	utils.WriteStandardJSON(w, r, http.StatusCreated, response)
}

// ListUsers handles user listing with pagination
func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "ListUsers request received")

	offset, limit := utils.ParsePagination(r)
	list, err := c.svc.ListUsers(ctx, offset, limit)
	if err != nil {
		code := errors.GetCode(err)
		if code == "" {
			code = constants.InternalError
		}
		logger.LogError(ctx, "failed to list users", err)
		utils.WriteStandardJSON(w, r, http.StatusInternalServerError, dto.ErrorResult{
			Msg: err.Error(),
		}, string(code))
		return
	}

	logger.LogInfo(ctx, "users listed successfully")
	response := dto.ToUserListResponse(list, len(list), offset, limit)
	utils.WriteStandardJSON(w, r, http.StatusOK, response)
}

// GetUser handles retrieving a single user by ID
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.LogInfo(ctx, "GetUser request received")

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		logger.LogWarn(ctx, "user id is required but not provided")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "user id is required",
		}, string(constants.InvalidParameter))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.LogWarn(ctx, "invalid user id format")
		utils.WriteStandardJSON(w, r, http.StatusBadRequest, dto.ErrorResult{
			Msg: "invalid user id format",
		}, string(constants.InvalidParameter))
		return
	}

	u, err := c.svc.GetUser(ctx, id)
	if err != nil {
		// Extract error code and determine HTTP status
		code := errors.GetCode(err)
		var httpStatus int

		switch code {
		case constants.NotFound:
			httpStatus = http.StatusNotFound
			logger.LogWarn(ctx, "user not found")
		default:
			httpStatus = http.StatusInternalServerError
			logger.LogError(ctx, "internal error in get user", err)
		}

		utils.WriteStandardJSON(w, r, httpStatus, dto.ErrorResult{
			Msg: err.Error(),
		}, string(code))
		return
	}

	logger.LogInfo(ctx, "user retrieved successfully")
	response := dto.ToUserResponse(u)
	utils.WriteStandardJSON(w, r, http.StatusOK, response)
}
