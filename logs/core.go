package logs

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"os"
)

// GetFileCore 获取文件日志core
func GetFileCore(encoder zapcore.EncoderConfig, logPath string) zapcore.Core {
	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	// 根据日志级别获取不同的core
	if l == zapcore.DebugLevel {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			zapcore.AddSync(fileWriter),
			l,
		)
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(fileWriter),
		l,
	)
}

// GetConsoleCore 获取控制台日志core
func GetConsoleCore(encoder zapcore.EncoderConfig) zapcore.Core {
	// 根据日志级别获取不同的编码器
	if l == zapcore.DebugLevel {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoder),
			zapcore.AddSync(os.Stdout),
			l,
		)
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		l,
	)
}
