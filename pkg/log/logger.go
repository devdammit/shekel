package log

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const envLogLevel = "LOG_LEVEL"

type Formatter string

const (
	FormatterText Formatter = "text"
	FormatterJSON Formatter = "json"
)

const (
	InfoLevel  = zap.InfoLevel  // 0, default level
	WarnLevel  = zap.WarnLevel  // 1
	ErrorLevel = zap.ErrorLevel // 2
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel // 5
	DebugLevel = zap.DebugLevel // -1
)

type Level = zapcore.Level

type Logger struct {
	l      *zap.Logger
	fields []Field
}

func NewLogger(o ...Option) Logger {
	var opts options
	for _, opt := range o {
		opt(&opts)
	}

	var encConfig zapcore.EncoderConfig
	if opts.Formatter == nil || *opts.Formatter == FormatterJSON {
		encConfig = zap.NewProductionEncoderConfig()
		encConfig.MessageKey = "message"
		encConfig.TimeKey = "@timestamp"
		encConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
		encConfig.CallerKey = "line_number"
		encConfig.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			if level == zapcore.WarnLevel {
				encoder.AppendString("warning")
				return
			}
			zapcore.LowercaseLevelEncoder(level, encoder)
		}
	} else {
		encConfig = zap.NewDevelopmentEncoderConfig()
		encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoder := zapcore.NewConsoleEncoder(encConfig)

	var writer io.Writer
	if opts.Writer != nil {
		writer = opts.Writer
	} else {
		writer = os.Stderr
	}

	var level Level
	if opts.Level != nil {
		level = *opts.Level
	} else if l, err := zapcore.ParseLevel(os.Getenv(envLogLevel)); err == nil {
		level = l
	} else {
		level = InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)
	zapLogger := zap.
		New(core).
		WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	return Logger{l: zapLogger}
}

func NewNilLogger() Logger {
	return Logger{l: zap.NewNop()}
}

type ctxFields struct{}

func NewContext(parent context.Context, fields ...Field) context.Context {
	val := parent.Value(ctxFields{})
	if val == nil {
		return context.WithValue(parent, ctxFields{}, fields)
	}
	return context.WithValue(parent, ctxFields{}, append(val.([]Field), fields...))
}

func (l Logger) WithContext(ctx context.Context) Logger {
	fields := ctx.Value(ctxFields{})
	if fields == nil {
		return l
	}
	newl := Logger{l: l.l, fields: append(fields.([]Field), l.fields...)}
	return newl
}

func (l Logger) With(fields ...Field) Logger {
	newl := Logger{l: l.l, fields: append(fields, l.fields...)}
	return newl
}

func (l Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, append(fields, l.fields...)...)
}
func (l Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, append(fields, l.fields...)...)
}
func (l Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, append(fields, l.fields...)...)
}
func (l Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, append(fields, l.fields...)...)
}
func (l Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, append(fields, l.fields...)...)
}
func (l Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, append(fields, l.fields...)...)
}
func (l Logger) Log(logLevel Level, msg string, fields ...Field) {
	l.l.Log(logLevel, msg, append(fields, l.fields...)...)
}

func (l Logger) ParseLevel(level string) (Level, error) {
	return zapcore.ParseLevel(level)
}
