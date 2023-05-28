package log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/wizke/go-util"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type (
	Fields   map[string]interface{}
	Level    int
	LevelStr string
	Mode     string
	Config   struct {
		LogLevel       Level
		LogFile        string
		SetColor       bool
		DayCount       int
		SessionKey     string
		Mode           Mode
		InstanceId     string
		InstanceIdShow bool
		Stdout         bool
	}
	CtxKey string
	Logger interface {
		Info(Fields)
		Warn(Fields)
		Error(Fields)
	}
	LoggerImpl struct {
		Mode Mode
	}
)

const (
	CommonMode Mode = ""
	JsonMode   Mode = "json"
)

const (
	CtxKeyLogWith = "log_with"
	CtxFields     = "fields"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

const (
	SqlLevel   Level = iota // Sql仅用于接管Sql日志输出
	GinLevel                // Start仅用于标记服务启动时的位置点与结束点
	StartLevel              // Start仅用于标记服务启动时的位置点与结束点
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// String
func (l Level) String() string {
	return [...]string{"SQL  ", "GIN  ", "Start", "Panic", "Fatal", "Error", "Warn ", "Info ", "Debug", "Trace"}[l]
}

func (l Level) StringLowerOnly() string {
	s := l.String()
	s = strings.ToLower(s)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func (l Level) EnumIndex() int {
	return int(l)
}

func NewLevel(str string) Level {
	p := strings.ToLower(str)
	p = strings.Replace(p, " ", "", -1)
	switch p {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return TraceLevel
	}
}

func NewMode(str string) Mode {
	p := strings.ToLower(str)
	p = strings.Replace(p, " ", "", -1)
	switch str {
	case "json":
		return JsonMode
	default:
		return CommonMode
	}
}

func (s LevelStr) GetLevel() Level {
	return NewLevel(string(s))
}

var ( //初始化修改后不再进行修改的全局参数
	instanceID     string //所在机器标识
	instanceIdShow bool   //是否输出机器标识
	isColor        bool
	isStdout       bool
	sessionKey     string
	logLevel       = TraceLevel
	logDaysCount   = 10
	logMode        = CommonMode
)

func SetLogLevel(l Level) {
	logLevel = l
}

func init() {
	instanceID, _ = os.Hostname()
	_ = InitLogger(Config{TraceLevel, "", isColor, logDaysCount, "", CommonMode, instanceID, false, false})
}

func InitLogger(config Config) error {
	if config.LogLevel == 0 {
		config.LogLevel = TraceLevel
	}
	if config.DayCount == 0 {
		config.DayCount = logDaysCount
	}
	if config.SessionKey != "" {
		if config.SessionKey == CtxKeyLogWith {
			return errors.New("session key(" + config.SessionKey + ") is not support, please input another one")
		}
		sessionKey = config.SessionKey
	}
	if config.Mode != "" {
		logMode = config.Mode
	}
	if config.InstanceId != "" {
		instanceID = config.InstanceId
	}

	instanceIdShow = config.InstanceIdShow
	logLevel = config.LogLevel
	logDaysCount = config.DayCount
	isColor = config.SetColor

	switch logMode {
	case CommonMode:
		log.SetFlags(log.Ldate | log.Lmicroseconds)
	case JsonMode:
		log.SetFlags(0)
	}
	switch runtime.GOOS {
	case "windows":
		// 默认在windows下启用开发模式,不进行日志文件写入，将日志输出到stdout
	case "linux", "darwin":
		// 在linux下启用标准输出验证如果开启则标准输出也会输出日志内容
		if config.Stdout {
			isStdout = config.Stdout
		}
		if config.LogFile != "" {
			path := config.LogFile + ".%Y%m%d"
			writer, _ := rotate.New(
				path,
				rotate.WithLinkName(config.LogFile),
				rotate.WithMaxAge(time.Duration(24*logDaysCount)*time.Hour),
				rotate.WithRotationTime(time.Duration(24)*time.Hour),
			)
			log.SetOutput(writer)
		}
	}
	return nil
}

func langFileStrToShortStr(fileStr string, maxLength int) (outStr string) {
	fileStrList := strings.Split(fileStr, "/")
	fileStrListLen := len(fileStrList)
	count := 0
	for index, str := range fileStrList {
		if index == fileStrListLen-1 {
			break
		}
		if str == "" {
			continue
		} else {
			if count == 0 {
				outStr += str + "/"
				count++
			} else {
				outStr += str[:1] + "/"
			}
		}
	}
	outStr += fileStrList[fileStrListLen-1]
	outStrLen := len(outStr)
	if outStrLen > maxLength {
		outStr = outStr[outStrLen-maxLength:]
	}
	return
}

func logCommon(level Level, file string, ctx context.Context, args ...interface{}) {
	ctxStr := ""
	ctxLogWithValue := ""
	if ctx != nil {
		sid := ctx.Value(sessionKey)
		if sid != nil {
			ctxStr = sid.(string)
		}
		logWithKey := ctx.Value(CtxKeyLogWith)
		if logWithKey != nil {
			ctxLogWithValue = logWithKey.(string)
		}
	}
	file = langFileStrToShortStr(file, 20)
	argsStr := ""
	switch ctxLogWithValue {
	case CtxFields:
		if len(args) > 0 {
			argsStr = fmt.Sprintf("%v", args[0])
		}
	default:
		for i, arg := range args {
			if i == 0 {
				argsStr += fmt.Sprintf("%v", arg)
			} else {
				argsStr += fmt.Sprintf(" %v", arg)
			}
		}
	}
	showStr := ""
	switch logMode {
	case JsonMode:
		levelStr := strings.ToLower(level.String())
		showStr = fmt.Sprintf(`{"time":"%s",level":"%s","file":"%s",%s%s"%s":%s}`,
			time.Now().Format("2006/01/02 15:04:05.000000"), levelStr, file,
			util.If(instanceIdShow, `"instance":"`+instanceID+`",`, ""),
			util.If(ctxStr != "", `"`+sessionKey+`":"`+ctxStr+`",`, ""),
			util.If(ctxLogWithValue == CtxFields, CtxFields, "msg"),
			util.If(ctxLogWithValue == CtxFields, argsStr, `"`+argsStr+`"`))
	default:
		instanceIdDisplay := util.If(instanceIdShow, "["+instanceID+"] ", "") // 是否要显示机器标识
		fileDisplay := util.If(isColor, fmt.Sprintf("%s%s%s", Cyan, file, Reset), file)
		levelStr := level.String()
		if isColor { // 在彩色输出模式下将pkg包中调用的日志也以非彩色形式输出
			color := ""
			switch level {
			case SqlLevel, GinLevel:
				color = Blue
			case StartLevel:
				color = Green
			case DebugLevel:
				color = Magenta
			case WarnLevel:
				color = Yellow
			case ErrorLevel:
				color = Red
			}
			levelStr = fmt.Sprintf("%s", util.If(isColor, fmt.Sprintf("%s%s%s", color, levelStr, Reset), levelStr))
		}
		showStr = fmt.Sprintf("%s%-20s [%s] %s%s", instanceIdDisplay, fileDisplay, levelStr,
			util.If(ctxStr != "", sessionKey+"="+ctxStr+" ", ""), argsStr)
	}
	log.Println(showStr)
	if isStdout {
		fmt.Println(showStr)
	}
}

func Starting(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logCommon(StartLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func Trace(args ...interface{}) {
	if logLevel < TraceLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(TraceLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func Printf(str string, args ...interface{}) {
	log.Printf(str, args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}

func Fatal(args ...interface{}) {
	if logLevel < FatalLevel {
		return
	}
	log.Fatal(args...)
}

func Panic(args ...interface{}) {
	if logLevel < PanicLevel {
		return
	}
	log.Panic(args...)
}

func SQL(format string, args []interface{}) {
	fileAndLine := fmt.Sprintf("%s", args[0])
	format = strings.Replace(format, "%s\n", "", 1)
	showStr := fmt.Sprintf(format, args[1:]...)
	logCommon(SqlLevel, fileAndLine, nil, showStr)
}

func GIN(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logCommon(GinLevel, fmt.Sprintf("%s:%d", file, line), nil, args...)
}

func WithFields(ctx context.Context, fields Fields, level Level) {
	if logLevel < level {
		return
	}
	_, file, line, _ := runtime.Caller(1)

	var fieldsJson string
	if fields != nil {
		fieldsJsonByte, err := json.Marshal(fields)
		if err != nil {
			fieldsJson = fmt.Sprintf("%v", fields)
		} else {
			fieldsJson = string(fieldsJsonByte)
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	logCommon(level, fmt.Sprintf("%s:%d", file, line), context.WithValue(ctx, CtxKeyLogWith, CtxFields), fieldsJson)
}
