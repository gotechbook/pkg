package logger

import (
	"context"
	"fmt"
	"os"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = "msg"

// Option is Helper option.
type Option func(*Helper)

// Helper is a logger helper.
type Helper struct {
	logger Logger
	msgKey string
}

// WithMessageKey with message key.
func WithMessageKey(k string) Option {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// NewHelper new a logger helper.
func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey: DefaultMessageKey, // default message key
		logger: logger,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey: h.msgKey,
		logger: WithContext(ctx, h.logger),
	}
}

// Log Print log by level and keyValue.
func (h *Helper) Log(level Level, keyValue ...interface{}) {
	_ = h.logger.Log(level, keyValue...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	_ = h.logger.Log(LevelDebug, h.msgKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelDebug, h.msgKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(keyValue ...interface{}) {
	_ = h.logger.Log(LevelDebug, keyValue...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	_ = h.logger.Log(LevelInfo, h.msgKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	_ = h.logger.Log(LevelInfo, h.msgKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(keyValue ...interface{}) {
	_ = h.logger.Log(LevelInfo, keyValue...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	_ = h.logger.Log(LevelWarn, h.msgKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelWarn, h.msgKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(keyValue ...interface{}) {
	_ = h.logger.Log(LevelWarn, keyValue...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	_ = h.logger.Log(LevelError, h.msgKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelError, h.msgKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(keyValue ...interface{}) {
	_ = h.logger.Log(LevelError, keyValue...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (h *Helper) Fatalw(keyValue ...interface{}) {
	_ = h.logger.Log(LevelFatal, keyValue...)
	os.Exit(1)
}
