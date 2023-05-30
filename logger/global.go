package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
)

var global = &loggerAppliance{}

type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global
}

// Log Print log by level and keyValue.
func Log(level Level, keyValue ...interface{}) {
	_ = global.Log(level, keyValue...)
}

// Context with context logger.
func Context(ctx context.Context) *Helper {
	return NewHelper(WithContext(ctx, global.Logger))
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyValue ...interface{}) {
	_ = global.Log(LevelDebug, keyValue...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyValue ...interface{}) {
	_ = global.Log(LevelInfo, keyValue...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyValue ...interface{}) {
	_ = global.Log(LevelWarn, keyValue...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyValue ...interface{}) {
	_ = global.Log(LevelError, keyValue...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyValue ...interface{}) {
	_ = global.Log(LevelFatal, keyValue...)
	os.Exit(1)
}

func init() {
	global.SetLogger(DefaultLogger)
}
