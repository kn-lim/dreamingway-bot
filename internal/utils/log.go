package utils

import (
	"go.uber.org/zap"
)

// Logger is the global logger for the application
var Logger *zap.SugaredLogger

// NewLogger creates a new logger
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	// nolint:errcheck
	defer logger.Sync()

	return logger.Sugar(), nil
}
