package logs

import (
	"context"
	"fmt"
	"github.com/oho-panda/utils/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *CustomLogger
	l      zapcore.Level
	s      = "default"
)

// ParseLevel 解析日志级别
func ParseLevel(level string) {
	parseLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		parseLevel = zapcore.InfoLevel
	}
	l = parseLevel
}

// CustomLogger 自定义日志
type CustomLogger struct {
	*zap.Logger
	gormLogger *GormLogger
	cronLogger *CronLogger
}

// GetLogger 获取日志
func GetLogger() *CustomLogger {
	return logger
}

// GetLevel 获取日志级别
func GetLevel() zapcore.Level {
	return l
}

// InitLogs 初始化日志
func InitLogs(service string, zapCore ...zapcore.Core) {
	if l == 0 {
		l = zapcore.InfoLevel
	}
	s = service
	cstLogger := &CustomLogger{
		gormLogger: &GormLogger{},
		cronLogger: &CronLogger{},
	}
	if len(zapCore) == 0 {
		consoleCore := GetConsoleCore(*GetConsoleEncoder())
		if consoleCore == nil {
			return
		}
		zapCore = append(zapCore, consoleCore)
	}
	core := zapcore.NewTee(
		zapCore...,
	)

	cstLogger.Logger = zap.New(
		core,
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	).With(zap.String("service", service))

	logger = cstLogger
}

// CtxLog 日志打印
func ctxLog(ctx context.Context, level zapcore.Level, msg string, v ...any) {
	if logger == nil {
		// 初始化日志
		InitLogs(s)
	}

	defer logger.Sync()
	// 获取traceId
	value := ctx.Value(consts.TraceIdKey)
	trace, ok := value.(string)
	var field []zap.Field
	if ok {
		field = append(field, zap.String(consts.TraceIdKey, trace))
	}
	sprintf := fmt.Sprintf(msg, v...)
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(sprintf, field...)
	case zapcore.InfoLevel:
		logger.Info(sprintf, field...)
	case zapcore.WarnLevel:
		logger.Warn(sprintf, field...)
	case zapcore.ErrorLevel:
		logger.Error(sprintf, field...)
	default:
		logger.Info(sprintf, field...)
	}
}

// CtxDebug 日志打印
func CtxDebug(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, zapcore.DebugLevel, msg, v...)
}

// CtxInfo 日志打印
func CtxInfo(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, zapcore.InfoLevel, msg, v...)
}

// CtxWarn 日志打印
func CtxWarn(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, zapcore.WarnLevel, msg, v...)
}

// CtxError 日志打印
func CtxError(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, zapcore.ErrorLevel, msg, v...)
}
