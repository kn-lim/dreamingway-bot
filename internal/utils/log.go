package utils

import (
	"go.uber.org/zap"
)

// Logger is the global logger for the application
var Logger *zap.SugaredLogger

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}
