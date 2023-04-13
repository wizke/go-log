package log

import (
	"fmt"
	"runtime"
)

func Warn(args ...interface{}) {
	if logLevel < WarnLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(WarnLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func WarnPrintf(str string, args ...interface{}) {
	if logLevel < WarnLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(WarnLevel, fmt.Sprintf("%s:%d", file, line), nil, fmt.Sprintf(str, args...))
}
