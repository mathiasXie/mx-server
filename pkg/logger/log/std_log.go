package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mathiasXie/gin-web/consts"
)

type STDLogger struct {
	mu      sync.Mutex
	level   Level
	logger  *log.Logger
	fields  map[string]interface{}
	context context.Context
}

// NewSTDLogger 创建一个新的标准日志实例
func NewSTDLogger() Logger {
	return &STDLogger{
		level:  Info,
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds),
		fields: make(map[string]interface{}),
	}
}

func (l *STDLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *STDLogger) log(level Level, ctx context.Context, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	levelStr := map[Level]string{
		Debug: "DEBUG",
		Info:  "INFO",
		Warn:  "WARN",
		Error: "ERROR",
		Fatal: "FATAL",
	}[level]

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	formattedMessage := fmt.Sprintf(format, args...)

	fieldsStr := ""
	for k, v := range l.fields {
		fieldsStr += fmt.Sprintf(" %s=%v", k, v)
	}

	contextStr := ""
	if ctx != nil {
		if reqID, ok := ctx.Value(consts.LogID).(string); ok {
			contextStr = fmt.Sprintf(" trace-id=%s", reqID)
		}
	}

	finalMessage := fmt.Sprintf("[%s] %s%s%s %s", timestamp, levelStr, contextStr, fieldsStr, formattedMessage)
	l.logger.Println(finalMessage)

	if level == Fatal {
		os.Exit(1)
	}
}

func (l *STDLogger) Debug(ctx context.Context, args ...interface{}) {
	l.log(Debug, ctx, fmt.Sprint(args...))
}

func (l *STDLogger) Info(ctx context.Context, args ...interface{}) {
	l.log(Info, ctx, fmt.Sprint(args...))
}

func (l *STDLogger) Warn(ctx context.Context, args ...interface{}) {
	l.log(Warn, ctx, fmt.Sprint(args...))
}

func (l *STDLogger) Error(ctx context.Context, args ...interface{}) {
	l.log(Error, ctx, fmt.Sprint(args...))
}

func (l *STDLogger) Panic(ctx context.Context, args ...interface{}) {
	l.log(Error, ctx, fmt.Sprint(args...))
}

func (l *STDLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log(Debug, ctx, format, args...)
}

func (l *STDLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log(Info, ctx, format, args...)
}

func (l *STDLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log(Warn, ctx, format, args...)
}

func (l *STDLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log(Error, ctx, format, args...)
}

func (l *STDLogger) WithField(key string, value interface{}) Logger {
	return l.WithFields(map[string]interface{}{key: value})
}
func (l *STDLogger) WithFields(fields map[string]interface{}) Logger {
	data := make(map[string]interface{}, len(l.fields)+len(fields))
	for k, v := range l.fields {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	l.fields = data
	return l
}

func (l *STDLogger) WithContext(ctx context.Context) Logger {
	newLogger := *l
	newLogger.context = ctx
	return &newLogger
}
