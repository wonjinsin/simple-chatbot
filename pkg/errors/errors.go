package errors

import (
	"errors"
	"fmt"

	pkgConstants "github.com/wonjinsin/simple-chatbot/internal/constants"
)

// CustomError represents an error with a 4-digit code
type CustomError struct {
	Code    pkgConstants.ErrorCode // 4-digit code (e.g., "0201")
	Message string
}

// Error implements error interface
func (e *CustomError) Error() string {
	return e.Message
}

// New creates a new CustomError with code and message.
// If an underlying error is provided, it combines the messages.
func New(code pkgConstants.ErrorCode, message string, err error) *CustomError {
	finalMessage := message
	if err != nil {
		finalMessage = fmt.Sprintf("%s: %s", message, err.Error())
	}
	return &CustomError{
		Code:    code,
		Message: finalMessage,
	}
}

// Wrap wraps an existing error with context.
// Accepts an optional error code. If provided, uses that code; otherwise preserves existing code or
// uses InternalError.
func Wrap(err error, message string, code ...pkgConstants.ErrorCode) error {
	if err == nil {
		return nil
	}

	// Determine which code to use
	var finalCode pkgConstants.ErrorCode
	if len(code) > 0 && code[0] != "" {
		// Use provided code
		finalCode = code[0]
	} else {
		// If already CustomError, preserve its code
		var customErr *CustomError
		if errors.As(err, &customErr) {
			finalCode = customErr.Code
		} else {
			// Otherwise use generic internal error
			finalCode = pkgConstants.InternalError
		}
	}

	return &CustomError{
		Code:    finalCode,
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
	}
}

// GetCode extracts the error code from CustomError
func GetCode(err error) pkgConstants.ErrorCode {
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return customErr.Code
	}
	return pkgConstants.UnknownError
}

// HasCode checks if error has specific code
func HasCode(err error, code pkgConstants.ErrorCode) bool {
	return GetCode(err) == code
}
