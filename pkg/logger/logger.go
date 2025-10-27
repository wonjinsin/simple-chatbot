package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/wonjinsin/go-boilerplate/pkg/constants"
)

// Initialize sets up the global logger
func Initialize(env string) {
	// Set time format to YYYY/MM/DD HH:MM:SS.mmm (e.g., 2025/01/01 01:01:01.333)
	zerolog.TimeFieldFormat = "2006/01/02 15:04:05.000"

	// JSON output for all environments
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Set log level based on environment
	if env == "local" || env == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// GetTrIDFromContext extracts TrID from context
func GetTrIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if trID, ok := ctx.Value(constants.ContextKeyTrID).(string); ok {
		return trID
	}
	return ""
}

// LogError logs an error with TrID from context
func LogError(ctx context.Context, msg string, err error) {
	trID := GetTrIDFromContext(ctx)
	if trID != "" {
		log.Error().
			Str("trid", trID).
			Err(err).
			Msg(msg)
	} else {
		log.Error().
			Err(err).
			Msg(msg)
	}
}

// LogInfo logs an info message with TrID from context
func LogInfo(ctx context.Context, msg string) {
	trID := GetTrIDFromContext(ctx)
	if trID != "" {
		log.Info().
			Str("trid", trID).
			Msg(msg)
	} else {
		log.Info().
			Msg(msg)
	}
}

// LogWarn logs a warning message with TrID from context
func LogWarn(ctx context.Context, msg string) {
	trID := GetTrIDFromContext(ctx)
	if trID != "" {
		log.Warn().
			Str("trid", trID).
			Msg(msg)
	} else {
		log.Warn().
			Msg(msg)
	}
}

// LogDebug logs a debug message with TrID from context
func LogDebug(ctx context.Context, msg string) {
	trID := GetTrIDFromContext(ctx)
	if trID != "" {
		log.Debug().
			Str("trid", trID).
			Msg(msg)
	} else {
		log.Debug().
			Msg(msg)
	}
}

// WithFields returns a logger with additional fields and TrID from context
func WithFields(ctx context.Context, fields map[string]interface{}) *zerolog.Event {
	trID := GetTrIDFromContext(ctx)
	event := log.Info()

	if trID != "" {
		event = event.Str("trid", trID)
	}

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	return event
}
