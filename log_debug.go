package log

import (
	"fmt"
	"runtime"
)

func Debug(args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(DebugLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func DebugPrintf(str string, args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(DebugLevel, fmt.Sprintf("%s:%d", file, line), nil, fmt.Sprintf(str, args...))
}
