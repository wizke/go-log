package log

import (
	"fmt"
	"runtime"
)

func Info(args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(InfoLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func InfoPrintf(str string, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(InfoLevel, fmt.Sprintf("%s:%d", file, line), nil, fmt.Sprintf(str, args...))
}
