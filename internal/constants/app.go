package constants

// Context keys (use unexported type for safety)
type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyTrID   contextKey = "tr_id"
)
