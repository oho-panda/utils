package logs

import (
	"context"
	"testing"
)

var (
	ctx context.Context
)

func init() {
	ctx = context.WithValue(context.Background(), "trace_id", "123321")
}

func TestNoInitLogs(t *testing.T) {
	CtxInfo(ctx, "测试哦哦哦哦哦")
}

func TestInitLogs(t *testing.T) {
	InitLogs("test_service")
	CtxInfo(ctx, "测试哦哦哦哦哦")
}

func TestErrorLogs(t *testing.T) {
	ParseLevel("error")
	InitLogs("TestErrorLogs")
	CtxInfo(ctx, "测试哦哦哦哦哦")
	CtxWarn(ctx, "测试哦哦哦哦哦")
	CtxError(ctx, "测试哦哦哦哦哦")
}

func TestDebugLogs(t *testing.T) {
	ParseLevel("info")
	logPath := "./l/test.log"
	encoder := GetFileEncoder()
	core := GetFileCore(*encoder, logPath)
	consoleEncoder := GetConsoleEncoder()
	consoleCore := GetConsoleCore(*consoleEncoder)
	InitLogs("TestDebugLogs", core, consoleCore)
	CtxInfo(ctx, "测试哦哦哦哦哦")
	CtxWarn(ctx, "测试哦哦哦哦哦")
	CtxError(ctx, "测试哦哦哦哦哦")
}

func TestProductLogs(t *testing.T) {
	ParseLevel("warn")
	logPath := "./l/test.log"
	encoder := GetFileEncoder()
	core := GetFileCore(*encoder, logPath)
	InitLogs("TestProductLogs", core)
	CtxInfo(ctx, "测试哦哦哦哦哦")
	CtxWarn(ctx, "测试哦哦哦哦哦")
	CtxError(ctx, "测试哦哦哦哦哦")
}
