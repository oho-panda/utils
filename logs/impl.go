package logs

import (
	"context"
	"go.uber.org/zap/zapcore"
	glog "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	LogLevel glog.LogLevel
}

type CronLogger struct {
}

func GLogger() *GormLogger {
	return logger.gormLogger
}

func CLogger() *CronLogger {
	return logger.cronLogger
}

// -------------- cronV3 interface实现 --------------

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	ctxLog(context.Background(), zapcore.InfoLevel, msg, keysAndValues...)
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, err.Error())
	ctxLog(context.Background(), zapcore.ErrorLevel, msg+" err %s", keysAndValues...)
}

// -------------- gorm interface实现 --------------

// LogMode 实现logger.Interface的LogMode方法
func (l *GormLogger) LogMode(level glog.LogLevel) glog.Interface {
	l.LogLevel = level
	return l
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, zapcore.InfoLevel, msg, data...)
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, zapcore.WarnLevel, msg, data...)
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, zapcore.ErrorLevel, msg, data...)
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		// 错误级别日志
		ctxLog(ctx, zapcore.ErrorLevel, "timeConsume: %s, rows: %d, sql: %s, error: %v", elapsed, rows, sql, err)
	} else {
		// 普通信息级别日志
		ctxLog(ctx, zapcore.InfoLevel, "timeConsume: %s, rows: %d, sql: %s", elapsed, rows, sql)
	}
}
