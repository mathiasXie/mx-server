package log

import (
	"context"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

type Logger interface {
	SetLevel(level Level)
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
}
