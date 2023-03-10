package gocom

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var loggerOnce sync.Once

func Logger() *zap.Logger {

	loggerOnce.Do(func() {

		conf := zap.NewProductionConfig()
		conf.EncoderConfig = zap.NewProductionEncoderConfig()
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, _ = conf.Build()
	})

	return logger
}
