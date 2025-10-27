package constants

type ErrorCode string

// Error codes - 4 digit format starting with 0, aligned with HTTP status codes
// 04xx: Client errors (matches HTTP 4xx)
// 05xx: Server errors (matches HTTP 5xx)

const (
	UnknownError ErrorCode = "0000" // HTTP 200 OK
	// Client errors (04xx)
	InvalidParameter ErrorCode = "0400" // HTTP 400 Bad Request
	NotFound         ErrorCode = "0404" // HTTP 404 Not Found
	ConstraintError  ErrorCode = "0409" // HTTP 409 Conflict

	// Server errors (05xx)
	InternalError ErrorCode = "0500" // HTTP 500 Internal Server Error
	DatabaseError ErrorCode = "0500" // HTTP 500 Internal Server Error
)
