package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(zapConfig zap.Config) error {
	config := zapConfig
	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	logger.Sync()
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
