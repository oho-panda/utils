package logs

import (
	"context"
	"fmt"
	"github.com/Free-D/utils/consts"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	glog "gorm.io/gorm/logger"
	"io"
	"os"
	"time"
)

var (
	logger *CustomLogger
)

type CustomLogger struct {
	*zap.Logger
	gormLogger *GormLogger
	cronLogger *CronLogger
}

type GormLogger struct {
	LogLevel glog.LogLevel
}

type CronLogger struct {
}

func GetLogger() *CustomLogger {
	return logger
}

func GLogger() *GormLogger {
	return logger.gormLogger
}

func CLogger() *CronLogger {
	return logger.cronLogger
}

type encoderType string

const (
	encoderTypeFile    encoderType = "file"
	encoderTypeConsole encoderType = "console"
)

func getSubEncoderByLevel(level zapcore.Level, t encoderType) (zapcore.LevelEncoder, zapcore.CallerEncoder) {
	if level == zapcore.DebugLevel {
		if t == encoderTypeFile {
			return zapcore.CapitalLevelEncoder, zapcore.FullCallerEncoder
		}
		return zapcore.CapitalColorLevelEncoder, zapcore.FullCallerEncoder
	}
	return zapcore.LowercaseLevelEncoder, zapcore.ShortCallerEncoder
}

func getEncoderByLevel(level zapcore.Level, t encoderType) zapcore.EncoderConfig {
	levelEncoder, callerEncoder := getSubEncoderByLevel(level, t)
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    levelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   callerEncoder,
	}
}

func getCoreByLevel(level zapcore.Level, encoder zapcore.EncoderConfig, writer io.Writer) zapcore.Core {
	return zapcore.NewCore(
		func(level zapcore.Level) zapcore.Encoder {
			if level == zapcore.DebugLevel {
				return zapcore.NewConsoleEncoder(encoder)
			}
			return zapcore.NewJSONEncoder(encoder)
		}(level),
		zapcore.AddSync(writer),
		level,
	)
}

// InitLogs 初始化日志
func InitLogs(prefix, path, app, level string) {
	cstLogger := &CustomLogger{
		gormLogger: &GormLogger{},
		cronLogger: &CronLogger{},
	}
	parseLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	logPath := path + "/" + app + ".log"

	fileEncoder := getEncoderByLevel(parseLevel, encoderTypeFile)
	consoleEncoder := getEncoderByLevel(parseLevel, encoderTypeConsole)

	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	fileCore := getCoreByLevel(parseLevel, fileEncoder, fileWriter)
	consoleCore := getCoreByLevel(parseLevel, consoleEncoder, os.Stdout)
	core := zapcore.NewTee(
		fileCore,
		consoleCore,
	)

	cstLogger.Logger = zap.New(
		core,
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	).With(zap.String("service", prefix+"_"+app))

	logger = cstLogger
}

func ctxLog(ctx context.Context, l *CustomLogger, level zapcore.Level, msg string, v ...any) {
	if l == nil {
		if logger == nil {
			// 初始化日志
			InitLogs("service", "./log", "debug", "debug")
		}
		l = logger
	}
	defer l.Sync()
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
		l.Debug(sprintf, field...)
	case zapcore.InfoLevel:
		l.Info(sprintf, field...)
	case zapcore.WarnLevel:
		l.Warn(sprintf, field...)
	case zapcore.ErrorLevel:
		l.Error(sprintf, field...)
	default:
		l.Info(sprintf, field...)
	}
}

func CtxDebug(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, logger, zapcore.DebugLevel, msg, v...)
}

func CtxInfo(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, logger, zapcore.InfoLevel, msg, v...)
}

func CtxWarn(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, logger, zapcore.WarnLevel, msg, v...)
}

func CtxError(ctx context.Context, msg string, v ...any) {
	ctxLog(ctx, logger, zapcore.ErrorLevel, msg, v...)
}

// -------------- cronV3 interface实现 --------------

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	ctxLog(context.Background(), logger, zapcore.InfoLevel, msg, keysAndValues...)
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, err.Error())
	ctxLog(context.Background(), logger, zapcore.ErrorLevel, msg+" err %s", keysAndValues...)
}

// -------------- gorm interface实现 --------------

// LogMode 实现logger.Interface的LogMode方法
func (l *GormLogger) LogMode(level glog.LogLevel) glog.Interface {
	l.LogLevel = level
	return l
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, logger, zapcore.InfoLevel, msg, data...)
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, logger, zapcore.WarnLevel, msg, data...)
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	ctxLog(ctx, logger, zapcore.ErrorLevel, msg, data...)
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		// 错误级别日志
		ctxLog(ctx, logger, zapcore.ErrorLevel, "timeConsume: %s, rows: %d, sql: %s, error: %v", elapsed, rows, sql, err)
	} else {
		// 普通信息级别日志
		ctxLog(ctx, logger, zapcore.InfoLevel, "timeConsume: %s, rows: %d, sql: %s", elapsed, rows, sql)
	}
}
