package log

import (
	"context"
	"fmt"
	"testing"
)

func commonLog() {
	Starting("start test")
	Trace("trace test")
	Debug("debug test")
	Info("info test")
	Warn("warn test")
	Error("error test")
	WithFields(nil, Fields{"1": "a"}, TraceLevel)
	WithFields(nil, Fields{"1": "a"}, DebugLevel)
	WithFields(nil, Fields{"1": "a"}, InfoLevel)
	WithFields(nil, Fields{"1": "a"}, WarnLevel)
	WithFields(nil, Fields{"1": "a"}, ErrorLevel)
	ctx := context.Background()
	ctx1 := context.WithValue(ctx, "session_id", "s1")
	InfoWithCtx(ctx1, "info")
	WarnWithCtx(ctx1, "warn")
	ErrorWithCtx(ctx1, "error")
	WithFields(ctx1, Fields{"1": "a", "2": "b"}, InfoLevel)
}

func TestLangFileStrToShortStr(t *testing.T) {
	fmt.Println(fmt.Sprintf("=====%q", langFileStrToShortStr("asdfasf", 20)))
}

func TestLogInit(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: InfoLevel,
		SetColor: true,
		DayCount: 10,
		LogFile:  "log/app.log",
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel:   NewLevel("DEBUG"),
		SetColor:   false,
		DayCount:   10,
		LogFile:    "log/app.log",
		SessionKey: "trace_id",
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

func TestColor(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: InfoLevel,
		SetColor: true,
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: InfoLevel,
		SetColor: false,
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}

func TestJsonMode(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		Mode:       JsonMode,
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("3>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   false,
		Mode:       NewMode("json"),
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("4<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
func TestLevel(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: TraceLevel,
		SetColor: true,
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: NewLevel("deBug"),
		SetColor: true,
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("3>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: NewLevel(" info"),
		SetColor: true,
	})
	commonLog()
	fmt.Println("3<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("4>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		LogLevel: NewLevel(" error"),
		SetColor: true,
	})
	commonLog()
	fmt.Println("4<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
func TestSessionKey(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "s_id",
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("3>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("3<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
func TestSetLogHandler(t *testing.T) {
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "session_id",
	})
	SetLogHandler(func(s string) string {
		return ">>> " + s + " <<<"
	})
	commonLog()
	fmt.Println("1<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "s_id",
	})
	commonLog()
	fmt.Println("2<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("3>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	_ = InitLogger(Config{
		SetColor:   true,
		SessionKey: "session_id",
	})
	commonLog()
	fmt.Println("3<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
