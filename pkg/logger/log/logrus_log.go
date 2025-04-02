package log

import (
	"context"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	entry   *logrus.Entry
	appName string
}

func NewLogrusLogger(conf *config.Config) Logger {
	err := os.MkdirAll(conf.Log.FileDirectory, os.ModePerm)
	if err != nil {
		panic(err)
		return nil
	}

	logFilePath := fmt.Sprintf("%s/%s-%%Y%%m%%d-%%H.log", conf.Log.FileDirectory, conf.AppName)
	writer, err := rotatelogs.New(
		logFilePath,
		rotatelogs.WithLinkName(fmt.Sprintf("%s/%s-latest.log", conf.Log.FileDirectory, conf.AppName)),
		rotatelogs.WithMaxAge(time.Duration(conf.Log.MaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
		return nil
	}

	logger := logrus.New()
	logger.SetOutput(writer)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // Custom timestamp format
	})
	logger.SetLevel(logrus.InfoLevel)
	return &LogrusLogger{entry: logrus.NewEntry(logger), appName: conf.AppName}
}

func (l *LogrusLogger) SetLevel(level Level) {
	l.entry.Logger.SetLevel(logrus.Level(level))
}

func (l *LogrusLogger) Debug(ctx context.Context, args ...interface{}) {
	l.withContext(ctx).entry.Debug(args...)
}

func (l *LogrusLogger) Info(ctx context.Context, args ...interface{}) {
	l.withContext(ctx).entry.Info(args...)
}

func (l *LogrusLogger) Warn(ctx context.Context, args ...interface{}) {
	l.withContext(ctx).entry.Warn(args...)
}

func (l *LogrusLogger) Error(ctx context.Context, args ...interface{}) {
	l.withContext(ctx).entry.Error(args...)
}

func (l *LogrusLogger) Panic(ctx context.Context, args ...interface{}) {
	l.withContext(ctx).entry.Panic(args...)
}

func (l *LogrusLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.withContext(ctx).entry.Debugf(format, args...)
}

func (l *LogrusLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.withContext(ctx).entry.Infof(format, args...)
}

func (l *LogrusLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.withContext(ctx).entry.Warnf(format, args...)
}

func (l *LogrusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.withContext(ctx).entry.Errorf(format, args...)
}

func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{entry: l.entry.WithField(key, value)}
}
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{entry: l.entry.WithFields(fields)}
}

func (l *LogrusLogger) WithContext(ctx context.Context) Logger {
	fields := logrus.Fields{}
	fields["trace_id"], _ = ctx.Value(consts.LogID).(string)
	fields["app_name"] = l.appName

	serverName, err := os.Hostname()
	if err != nil {
		serverName = "unknown"
	}
	fields["server_name"] = serverName

	return &LogrusLogger{entry: l.entry.WithFields(fields)}
}

// withContext 提取上下文字段并返回新的 Logger 实例
func (l *LogrusLogger) withContext(ctx context.Context) *LogrusLogger {
	if ctx == nil {
		return l
	}
	return l.WithContext(ctx).(*LogrusLogger)
}
