package worker

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msgf("%v", args...)
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

// Debug logs a message at Debug level.
func (logger *Logger) Debug(args ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf("%v", args...)
}

// Info logs a message at Info level.
func (logger *Logger) Info(args ...interface{}) {
	log.WithLevel(zerolog.InfoLevel).Msgf("%v", args...)
}

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...interface{}) {
	log.WithLevel(zerolog.WarnLevel).Msgf("%v", args...)
}

// Error logs a message at Error level.
func (logger *Logger) Error(args ...interface{}) {
	log.WithLevel(zerolog.ErrorLevel).Msgf("%v", args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	log.WithLevel(zerolog.FatalLevel).Msgf("%v", args...)
}
