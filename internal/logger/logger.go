package logger

import (
	"go.uber.org/zap"
)

// NewLogger creates and returns a structured Zap logger instance
// Zap is used for efficient, structured JSON logging in production
// Returns: *zap.Logger for high-performance logging throughout the application
func NewLogger() (*zap.Logger, error) {
	// Production config creates optimized logger for production environments
	// Outputs JSON format suitable for log aggregation systems (ELK, Splunk, etc.)
	return zap.NewProduction()
}
