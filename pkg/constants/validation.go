package constants

// Validation limits
const (
	// User validation
	MaxNameLength  = 200
	MinEmailLength = 3
	MaxEmailLength = 320 // RFC 5321 limit
)

// Regex patterns
const (
	EmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// ID generation
const (
	IDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)
