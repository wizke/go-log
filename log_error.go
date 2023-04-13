package log

import (
	"fmt"
	"runtime"
)

func Error(args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(ErrorLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func ErrorPrintf(str string, args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(ErrorLevel, fmt.Sprintf("%s:%d", file, line), nil, fmt.Sprintf(str, args...))
}
