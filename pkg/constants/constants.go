package constants

// Context keys (use unexported type for safety)
type ContextKey string

const (
	// ContextKeyTrID is the key for storing TrID in context
	ContextKeyTrID ContextKey = "tr_id"
)
