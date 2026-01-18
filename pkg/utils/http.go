package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/wonjinsin/simple-chatbot/pkg/constants"
)

// ParsePagination extracts pagination parameters from HTTP request
func ParsePagination(r *http.Request) (offset, limit int) {
	q := r.URL.Query()
	offset, _ = strconv.Atoi(q.Get("offset"))
	limit, _ = strconv.Atoi(q.Get("limit"))

	if offset < 0 {
		offset = constants.DefaultOffset
	}
	if limit <= 0 || limit > constants.MaxLimit {
		limit = constants.DefaultLimit
	}

	return offset, limit
}

// WriteJSON writes JSON response with proper headers
func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSONCharset)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

// ParseJSONBody parses JSON request body into the provided struct
func ParseJSONBody(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ExtractPathParam extracts path parameter from URL
// Example: ExtractPathParam("/users/123", "/users/") returns "123"
func ExtractPathParam(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}

// StandardResponse represents the standard HTTP response format
type StandardResponse struct {
	TrID   string `json:"trid"`
	Code   string `json:"code"`
	Result any    `json:"result,omitempty"`
}

// WriteStandardJSON writes a standard JSON response with TrID
// Accepts an optional custom code string. If not provided, uses HTTP status code.
func WriteStandardJSON(
	w http.ResponseWriter,
	r *http.Request,
	httpStatus int,
	result any,
	customCode ...string,
) {
	// Get TrID from context
	trID := ""
	if ctx := r.Context(); ctx != nil {
		if id, ok := ctx.Value(constants.ContextKeyTrID).(string); ok {
			trID = id
		}
	}

	// Determine response code: use custom code if provided, otherwise format HTTP status
	var codeStr string
	if len(customCode) > 0 && customCode[0] != "" {
		codeStr = customCode[0]
	} else {
		codeStr = fmt.Sprintf("%04d", httpStatus)
	}

	// Create standard response using struct (preserves field order)
	response := StandardResponse{
		TrID:   trID,
		Code:   codeStr,
		Result: result,
	}

	w.Header().Set(constants.HeaderContentType, constants.ContentTypeJSONCharset)
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
