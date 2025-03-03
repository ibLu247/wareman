package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func NewLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}
