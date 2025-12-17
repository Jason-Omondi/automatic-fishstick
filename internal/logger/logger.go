package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and returns a structured Zap logger instance
// In development: pretty-printed colored output to console
// In production: JSON format suitable for log aggregation systems
// Returns: *zap.Logger for high-performance logging throughout the application
func NewLogger() (*zap.Logger, error) {
	// Check environment to determine logging mode
	env := os.Getenv("ENV")
	if env == "production" {
		return newProductionLogger()
	}
	return newDevelopmentLogger()
}

// newDevelopmentLogger creates a pretty-printed logger for development
// Features: colored output, human-readable format, detailed stack traces
func newDevelopmentLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	// Custom encoder configuration for prettier output
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.ConsoleSeparator = " | "

	// Output to stdout
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	return config.Build()
}

// newProductionLogger creates a JSON logger for production
// Features: JSON format for log aggregation, optimized performance
func newProductionLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	// JSON encoder for structured logging
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	// Output to stdout
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	return config.Build()
}
