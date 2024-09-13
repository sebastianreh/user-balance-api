package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, keyValuePairs ...interface{})
	Warn(msg string, keyValuePairs ...interface{})
	Error(msg string, keyValuePairs ...interface{})
	ErrorAt(err error, layer, origin string)
	Fatal(msg string)
}

type logger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewLogger() Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	config.EncoderConfig.CallerKey = ""
	config.DisableStacktrace = true
	log, _ := config.Build()
	sugaredLogger := log.Sugar()

	return &logger{sugaredLogger: sugaredLogger}
}

func (l *logger) Info(msg string, keyValuePairs ...interface{}) {
	l.sugaredLogger.Infof(msg+" %s ", keyValuePairs...)
}

func (l *logger) Warn(msg string, keyValuePairs ...interface{}) {
	l.sugaredLogger.Warnf(msg+". origin: %s", keyValuePairs...)
}

func (l *logger) Error(msg string, keyValuePairs ...interface{}) {
	l.sugaredLogger.Errorf(msg+". origin: %s", keyValuePairs...)
}

func (l *logger) ErrorAt(err error, layer, origin string) {
	l.sugaredLogger.Errorf(err.Error()+". origin: %s", fmt.Sprintf("%s.%s", layer, origin))
}

func (l *logger) Fatal(msg string) {
	l.sugaredLogger.Fatalf(msg)
}
