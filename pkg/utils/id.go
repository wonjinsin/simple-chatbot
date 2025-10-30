package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/wonjinsin/simple-chatbot/pkg/constants"
)

// GenerateID generates a unique ID using timestamp and counter
func GenerateID(counter int64) string {
	n := time.Now().UnixNano()
	return FormatID(n, counter)
}

// FormatID formats timestamp and counter into a readable ID
func FormatID(timestamp int64, counter int64) string {
	var buf [32]byte
	i := len(buf)
	x := uint64((timestamp << 13) ^ (timestamp >> 7) ^ counter)

	for x > 0 {
		i--
		buf[i] = constants.IDAlphabet[x%36]
		x /= 36
	}
	return string(buf[i:])
}

// GenerateRandomID generates a cryptographically secure random ID
func GenerateRandomID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
