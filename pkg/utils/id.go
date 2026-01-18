package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	pkgConstants "github.com/wonjinsin/simple-chatbot/pkg/constants"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
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
		buf[i] = pkgConstants.IDAlphabet[x%36]
		x /= 36
	}
	return string(buf[i:])
}

// GenerateRandomID generates a cryptographically secure random ID
func GenerateRandomID(length int) (string, error) {
	if length <= 0 {
		return "", errors.New(constants.InvalidParameter, "length must be positive", nil)
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.Wrap(err, "failed to generate random bytes")
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
