package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Setup zap logger to be used as the logging engine,
//
//	for the various level of logging activities to be performed
func init() {
	var err error

	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig

	log, err = config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func Info(message string, args ...interface{}) {
	log.Info(fmt.Sprintf(message, args...))
}

func Error(message interface{}) {
	log.Error(fmt.Sprintf("%v", message))
}

func Debug(message string, args ...interface{}) {
	log.Debug(fmt.Sprintf(message, args...))
}

func Fatal(message interface{}) {
	log.Fatal(fmt.Sprintf("%v", message))
}
