package slogr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// HandlerTepe defienes the type of handler to use
type HandlerType int

const (
	HandlerTypeJSON HandlerType = iota
	HandlerTypeText
)

type Options struct {
	// Level threshold for logging
	Level slog.Level

	// Whether to prefix messages with level
	AddLevelPrefix bool

	// Type of handler to use
	HandlerType HandlerType

	// Custom handler if provided will override HandlerType
	CustomHandler slog.Handler

	// Additional handler options
	HandlerOptions *slog.HandlerOptions
}

func DefaultOptions() *Options {
	return &Options{
		Level:          slog.LevelInfo,
		AddLevelPrefix: false,
		HandlerType:    HandlerTypeText,
		HandlerOptions: &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				return a
			},
		},
	}
}

var defaultLogger = New(os.Stdout, DefaultOptions())

// SetOutput sets the output for the default logger
func SetOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

// SetLevel sets the level for the default logger
func SetLevel(level slog.Level) {
	defaultLogger.SetLevel(level)
}

func GetLevel() slog.Level {
	return defaultLogger.GetLevel()
}

type Logger struct {
	// level is the minimum level output by this Logger
	level slog.Level

	// shouldPreficMessageWithLevel is whether to include the log level prefix in each log.
	shouldPreficMessageWithLevel bool

	// handlerType is the type of handler used by this Logger
	handlerType HandlerType

	writerType io.Writer

	slogger *slog.Logger
}

// New create a new logger with the given threshold and output
func New(output io.Writer, opts *Options) *Logger {
	if opts == nil {
		opts = DefaultOptions()
	}

	if opts.HandlerOptions == nil {
		opts.HandlerOptions = &slog.HandlerOptions{
			Level: opts.Level,
		}
	}

	var handler slog.Handler

	// if custom handler is provided use it
	if opts.CustomHandler != nil {
		handler = opts.CustomHandler
	} else {
		switch opts.HandlerType {
		case HandlerTypeJSON:
			handler = slog.NewJSONHandler(output, opts.HandlerOptions)
		default:
			handler = slog.NewTextHandler(output, opts.HandlerOptions)
		}
	}

	if opts.AddLevelPrefix {
		handler = &levelPrefixHandler{handler}
	}

	return &Logger{
		level:                        opts.Level,
		shouldPreficMessageWithLevel: opts.AddLevelPrefix,
		handlerType:                  opts.HandlerType,
		slogger:                      slog.New(handler),
		writerType:                   output,
	}
}

// levelPrefixHandler adds a level prefix to log message
type levelPrefixHandler struct {
	slog.Handler
}

func (h *levelPrefixHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Message = "- " + r.Level.String() + " - " + r.Message
	return h.Handler.Handle(ctx, r)
}

func (logger *Logger) SetHandler(output io.Writer, handlerType HandlerType, opts *slog.HandlerOptions) {
	if opts == nil {
		opts = &slog.HandlerOptions{
			Level: logger.level,
		}
	}

	var handler slog.Handler
	switch handlerType {
	case HandlerTypeJSON:
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	if logger.shouldPreficMessageWithLevel {
		handler = &levelPrefixHandler{handler}
	}

	logger.slogger = slog.New(handler)
}

// SetOutput changes the output destination for the logger
func (logger *Logger) SetOutput(output io.Writer) {
	logger.SetHandler(output, logger.handlerType, &slog.HandlerOptions{
		Level: logger.level,
	})
}

func (logger *Logger) SetLevel(level slog.Level) {
	logger.level = level
	logger.SetHandler(logger.writerType, logger.handlerType, &slog.HandlerOptions{
		Level: level,
	})
}

func (logger *Logger) GetLevel() slog.Level {
	return logger.level
}

func (logger *Logger) GetHandlerType() HandlerType {
	return logger.handlerType
}

// SetCustomHandler allows setting a custom handler
func (logger *Logger) SetCustomHandler(handler slog.Handler) {
	if logger.shouldPreficMessageWithLevel {
		handler = &levelPrefixHandler{handler}
	}

	logger.slogger = slog.New(handler)
}

func (logger *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.slogger.Log(ctx, level, msg, args...)
}

func (logger *Logger) Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	// Format the message with fmt.Sprintf before passing to the logger
	formattedMsg := fmt.Sprintf(format, args...)
	logger.slogger.Log(ctx, level, formattedMsg)
}

func (logger *Logger) Debug(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelDebug, msg, args...)
}

func (logger *Logger) Debugf(ctx context.Context, format string, args ...any) {
	logger.Logf(ctx, slog.LevelDebug, format, args...)
}

func (logger *Logger) Info(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelInfo, msg, args...)
}

func (logger *Logger) Infof(ctx context.Context, format string, args ...any) {
	logger.Logf(ctx, slog.LevelInfo, format, args...)
}

func (logger *Logger) Warn(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelWarn, msg, args...)
}

func (logger *Logger) Warnf(ctx context.Context, format string, args ...any) {
	logger.Logf(ctx, slog.LevelWarn, format, args...)
}

func (logger *Logger) Error(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelError, msg, args...)
}

func (logger *Logger) Errorf(ctx context.Context, format string, args ...any) {
	logger.Logf(ctx, slog.LevelError, format, args...)
}

func (logger *Logger) Fatal(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelError+4, msg, args...)
}

func (logger *Logger) Fatalf(ctx context.Context, format string, args ...any) {
	logger.Logf(ctx, slog.LevelError+4, format, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	defaultLogger.Log(ctx, level, msg, args...)
}

func Logf(ctx context.Context, level slog.Level, format string, args ...any) {
	defaultLogger.Logf(ctx, level, format, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	defaultLogger.Debug(ctx, msg, args...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	defaultLogger.Debugf(ctx, format, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	defaultLogger.Info(ctx, msg, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	defaultLogger.Infof(ctx, format, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	defaultLogger.Warn(ctx, msg, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	defaultLogger.Warnf(ctx, format, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	defaultLogger.Error(ctx, msg, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	defaultLogger.Errorf(ctx, format, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	defaultLogger.Fatal(ctx, msg, args...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	defaultLogger.Fatalf(ctx, format, args...)
}
