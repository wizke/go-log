package log

import (
	"context"
	"fmt"
	"runtime"
)

func DebugWithCtx(ctx context.Context, args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(DebugLevel, fmt.Sprintf("%s:%d", file, line), ctx, args...)
}

func InfoWithCtx(ctx context.Context, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(InfoLevel, fmt.Sprintf("%s:%d", file, line), ctx, args...)
}

func WarnWithCtx(ctx context.Context, args ...interface{}) {
	if logLevel < WarnLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(WarnLevel, fmt.Sprintf("%s:%d", file, line), ctx, args...)
}

func ErrorWithCtx(ctx context.Context, args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(ErrorLevel, fmt.Sprintf("%s:%d", file, line), ctx, args...)
}
