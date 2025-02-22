package logr

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)


func TestNew(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
		wantErr bool
	}{
		{
			name: "nil options should use defaults",
			opts: nil,
			wantErr: false,
		},
		{
			name: "custom options",
			opts: &Options{
				Level: slog.LevelDebug,
				AddLevelPrefix: true,
				HandlerType: HandlerTypeJSON,
			},
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := New(buf, test.opts)
			if logger == nil {
				t.Errorf("expected logger to be non-nil")
			}
		})
	}
}

func TestLogger_Levels(t *testing.T) {
	tests := []struct {
		name string
		level slog.Level
		message string
		contains string
	}{
		{
			name: "debug level",
			level: slog.LevelDebug,
			message: "debug message",
			contains: "DEBUG",
		},
		{
			name: "info level",
			level: slog.LevelInfo,
			message: "info message",
			contains: "INFO",
		},
		{
			name: "warn level",
			level: slog.LevelWarn,
			message: "warn message",
			contains: "WARN",
		},
		{
			name: "error level",
			level: slog.LevelError,
			message: "error message",
			contains: "ERROR",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := New(buf, &Options{
				Level: test.level,
				AddLevelPrefix: true,
				HandlerType: HandlerTypeText,
			})
			logger.Log(context.Background(), test.level, test.message)
			output := buf.String()
			if !strings.Contains(output, test.contains) {
				t.Errorf("expected output to contain %s, got %s", test.contains, output)
			}
		})
	}
}


func TestLogger_HandlerType(t *testing.T) {
	tests := []struct {
		name string
		handlerType HandlerType
		checkJSON bool
	}{
		{
			name: "text handler",
			handlerType: HandlerTypeText,
			checkJSON: false,
		},
		{
			name: "json handler",
			handlerType: HandlerTypeJSON,
			checkJSON: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			logger := New(buf, &Options{			
				HandlerType: test.handlerType,
			})
			logger.Log(context.Background(), slog.LevelInfo, "test message")
			output := buf.String()
			if test.checkJSON {
				var jsonMap map[string]interface{}
				if err := json.Unmarshal([]byte(output), &jsonMap); err != nil {
					t.Errorf("expected output to be valid JSON, got %s", output)
				}
			} else {
				if !strings.Contains(output, "INFO") {
					t.Errorf("expected output to contain INFO, got %s", output)
				}
			}
		})
	}
}

func TestLogger_CustomHandler(t *testing.T) {
	buf := new(bytes.Buffer)
	customHandler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := New(buf, nil)
	logger.SetCustomHandler(customHandler)

	logger.Info(context.Background(), "test message")

	if buf.Len() == 0 {
		t.Errorf("expected output to be non-empty")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := New(buf, nil)

	// set level to error and try to log Info
	logger.SetLevel(slog.LevelError)
	logger.Infof(context.Background(), "this should not be logged")

	if buf.Len() > 0 {
		t.Error("expected no output for Info level message when level is set to Error")
	}

	logger.Errorf(context.Background(), "this should be logged")
	if buf.Len() == 0 {
		t.Error("expected output for Error level message")
	}
}


func TestLogger_SetThreshold(t *testing.T) {
	ctx := context.Background()
	buf := new(bytes.Buffer)
	logger := New(buf, nil)
	logger.SetLevel(slog.LevelError)
	logger.Debug(ctx, "this is debug")
	logger.Debugf(ctx, "this is debug %d", 1)
	logger.Info(ctx, "this is info")
	logger.Infof(ctx, "this is info %d", 1)
	logger.Warn(ctx, "this is warn")
	logger.Warnf(ctx, "this is warn %d", 1)
	// logger.Errorf(ctx, "this is error %d", 1)
	logger.Error(ctx, "this is error")
	output := buf.String()
	t.Logf("output: %s", output)
	if numberOfNewLines := strings.Count(output, "\n"); numberOfNewLines != 1 {
		t.Error("expected 1 newline, got", numberOfNewLines)
	}
	// level := logger.GetLevel()
	// if level != slog.LevelError {
	// 	t.Errorf("expected level to be %s, got %s", slog.LevelError, level)
	// }
	// handlerType := logger.GetHandlerType()
	// if handlerType != HandlerTypeText {
	// 	t.Errorf("expected handler type to be HandlertText, got %T", handlerType)
	// }
}

func TestLogger_SetOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	
	logger := New(buf1, nil)
	logger.Info(context.Background(), "first buffer")
	
	if buf1.Len() == 0 {
		t.Error("expected output in first buffer")
	}

	logger.SetOutput(buf2)
	logger.Info(context.Background(), "second buffer")
	
	initialBuf1Size := buf1.Len()
	if buf2.Len() == 0 {
		t.Error("expected output in second buffer")
	}
	if buf1.Len() != initialBuf1Size {
		t.Error("expected first buffer to remain unchanged")
	}
}