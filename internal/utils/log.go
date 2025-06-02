package utils

import (
	"go.uber.org/zap"
)

// Logger is the global logger for the application
var Logger *zap.SugaredLogger

// NewLogger creates a new logger
func NewLogger(verbose bool) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	if verbose {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	// nolint:errcheck
	defer logger.Sync()

	return logger.Sugar(), nil
}
