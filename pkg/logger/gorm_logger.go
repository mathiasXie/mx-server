package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		CtxInfo(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		CtxWarn(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		CtxError(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]interface{}{
		"sql":           sql,
		"rows_affected": rows,
		"elapsed_time":  elapsed.String(),
	}

	if err != nil {
		fields["error"] = err
		CtxError(ctx, "SQL Error", fields)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		CtxWarn(ctx, "SQL Slow Query", fields)
		return
	}

	CtxInfo(ctx, "SQL Query", fields)
}
