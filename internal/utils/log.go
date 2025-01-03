package utils

import (
	"log"

	"go.uber.org/zap"
)

// Logger is the global logger for the application
var Logger *zap.SugaredLogger

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := Logger.Sync(); err != nil {
			log.Printf("error syncing logger: %v", err)
		}
	}()

	return logger.Sugar(), nil
}
