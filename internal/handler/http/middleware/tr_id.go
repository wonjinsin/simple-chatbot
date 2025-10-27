package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/wonjinsin/go-boilerplate/pkg/constants"
)

// TrID returns a middleware that generates transaction IDs
// Format: YYYYMMDDHHMMSSmmm (date+time+milliseconds) + 5-digit random number
// Example: 2025010101010199912345 (23 digits total)
func TrID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate TrID
			trID := GenerateTrID()

			// Add to context
			ctx := context.WithValue(r.Context(), constants.ContextKeyTrID, trID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GenerateTrID generates a transaction ID
// Format: YYYYMMDDHHMMSSmmm + 5-digit random number
func GenerateTrID() string {
	now := time.Now()

	// Format: YYYYMMDDHHMMSSmmm (17 digits)
	timestamp := now.Format("20060102150405")
	milliseconds := fmt.Sprintf("%03d", now.Nanosecond()/1000000)

	// Generate 5-digit random number (00000-99999)
	randomNum := generateRandomNumber(5)

	return timestamp + milliseconds + randomNum
}

// generateRandomNumber generates a random number string with specified length
func generateRandomNumber(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based random if crypto fails
		return fmt.Sprintf("%0*d", length, time.Now().UnixNano()%powerOf10(length))
	}

	// Convert bytes to number string
	var num int64
	for _, b := range bytes {
		num = (num*256 + int64(b)) % powerOf10(length)
	}

	return fmt.Sprintf("%0*d", length, num)
}

// powerOf10 calculates 10^n
func powerOf10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// GetTrID extracts TrID from context
func GetTrID(ctx context.Context) string {
	if trID, ok := ctx.Value(constants.ContextKeyTrID).(string); ok {
		return trID
	}
	return ""
}
