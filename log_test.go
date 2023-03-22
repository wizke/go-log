package log

import "testing"

func TestLog(t *testing.T) {
	InitLogger(Config{
		LogLevel: InfoLevel,
		SetColor: true,
		DayCount: 10,
		LogFile:  "log/app.log",
	})
	Starting("Start")
	Trace("trace")
	Debug("zero")
	Info("first")
	Warn("two")
	Error("three")
	WithFields("ab", Fields{"ab": "cd"}, TraceLevel)
	WithFields("ab", Fields{"ab": "cd"}, DebugLevel)
	WithFields("ab", Fields{"ab": "cd"}, InfoLevel)
	WithFields("ab", Fields{"ab": "cd"}, WarnLevel)
	InitLogger(Config{
		LogLevel: NewLevel("DEBUG"),
		SetColor: true,
		DayCount: 10,
		LogFile:  "log/app.log",
	})
	Starting("Start New")
	Trace("trace")
	Debug("zero")
	Info("first")
	Warn("two")
	Error("three")
	WithFields("ab", Fields{"ab": "cd"}, TraceLevel)
	WithFields("ab", Fields{"ab": "cd"}, DebugLevel)
	WithFields("ab", Fields{"ab": "cd"}, InfoLevel)
	WithFields("ab", Fields{"ab": "cd"}, WarnLevel)
}
