package log

import (
	"context"
)

var defaultLogger = NewLogger()
var nilLogger = NewNilLogger()

func Default() Logger {
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
}

func Nil() Logger {
	return nilLogger
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}

func With(fields ...Field) Logger {
	return defaultLogger.With(fields...)
}

func Debug(msg string, fields ...Field) {
	defaultLogger.l.Debug(msg, append(fields, defaultLogger.fields...)...)
}
func Info(msg string, fields ...Field) {
	defaultLogger.l.Info(msg, append(fields, defaultLogger.fields...)...)
}
func Warn(msg string, fields ...Field) {
	defaultLogger.l.Warn(msg, append(fields, defaultLogger.fields...)...)
}
func Error(msg string, fields ...Field) {
	defaultLogger.l.Error(msg, append(fields, defaultLogger.fields...)...)
}
func Fatal(msg string, fields ...Field) {
	defaultLogger.l.Fatal(msg, append(fields, defaultLogger.fields...)...)
}
func Panic(msg string, fields ...Field) {
	defaultLogger.l.Panic(msg, append(fields, defaultLogger.fields...)...)
}
