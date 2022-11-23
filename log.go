package log

import (
	"encoding/json"
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
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (l Level) String() string {
	return [...]string{"Panic", "Fatal", "Error", "Warn", "Info", "Debug", "Trace"}[l]
}

func (l Level) EnumIndex() int {
	return int(l)
}

func (ls LevelStr) GetLevel() Level {
	switch ls {
	case "Panic":
		return PanicLevel
	case "Fatal":
		return FatalLevel
	case "Error":
		return ErrorLevel
	case "Warn":
		return WarnLevel
	case "Info":
		return InfoLevel
	case "Debug":
		return DebugLevel
	default:
		return TraceLevel
	}
}

var ( //初始化修改后不再进行修改的全局参数
	_thisFile     = "log.go"
	defaultLogger = ""   //缺省logger 名称
	instanceID    string //所在机器标识
	logMainPrefix = ""   // 用于显示日志输出文件路径，清除源码内路径前部内容
	logPkgPrefix  = ""   // 用于显示日志输出文件路径，清除pkg包中调用日志的文件路径前部内容
	logLevel      = TraceLevel
	isColor       = false
	logDaysCount  = 10
)

func init() {
	instanceID, _ = os.Hostname()
	InitLogger(TraceLevel, _thisFile, "", true, 10)
}

func InitLogger(setLogLevel Level, thisFile string, logFile string, setColor bool, daysCount int) {
	logLevel = setLogLevel
	logDaysCount = daysCount

	_, file, _, _ := runtime.Caller(1)
	if thisFile == _thisFile {
		logPkgPrefix = strings.Replace(file, thisFile, "", 1)
	} else {
		logMainPrefix = strings.Replace(file, thisFile, "", 1)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	switch runtime.GOOS {
	case "windows", "darwin":
		// 默认在windows和macos下启用开发模式，输出所有日志等级
		// 不进行日志文件写入，将日志输出到stdout，并开启彩色输出
		isColor = true
		return
	case "linux":
		isColor = setColor
	}
	if logFile != "" {
		_, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		path := logFile + ".%Y%m%d"
		writer, _ := rotate.New(
			path,
			rotate.WithLinkName(logFile),
			rotate.WithMaxAge(time.Duration(24*logDaysCount)*time.Hour),
			rotate.WithRotationTime(time.Duration(24)*time.Hour),
		)

		if err != nil {
			log.Printf("open log file error : %s", err.Error())
			return
		}
		log.SetOutput(writer)
	}
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

func logCommon(levelStr, file string, args ...interface{}) {
	isPkgLog := false
	if strings.Index(file, logPkgPrefix) == -1 {
		file = strings.Replace(file, logMainPrefix, "", 1)
	} else {
		isPkgLog = true
		file = strings.Replace(file, logPkgPrefix, "daisy/", 1)
	}
	showStr := ""
	for _, arg := range args {
		showStr += fmt.Sprintf("%v ", arg)
	}
	if isColor && !isPkgLog { // 在彩色输出模式下将pkg包中调用的日志也以非彩色形式输出
		log.Println(fmt.Sprintf("%s%-20s%s [%s]", Cyan, langFileStrToShortStr(file, 20), Reset, levelStr), showStr)
	} else {
		log.Println(fmt.Sprintf("%-20s [%s]", langFileStrToShortStr(file, 20), levelStr), showStr)
	}
}

func Starting(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logCommon(fmt.Sprintf("%s", util.If(isColor, Green, ""))+"Start"+Reset, fmt.Sprintf("%s:%d", file, line), args...)
}

func Trace(args ...interface{}) {
	if logLevel < TraceLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon("Trace", fmt.Sprintf("%s:%d", file, line), args...)
}

func Debug(args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon("Debug", fmt.Sprintf("%s:%d", file, line), args...)
}

func InfoForSQL(format string, args []interface{}) {
	fileAndLine := fmt.Sprintf("%s", args[0])
	format = strings.Replace(format, "%s\n", "", 1)
	showStr := fmt.Sprintf(format, args[1:]...)
	logCommon(fmt.Sprintf("%s", util.If(isColor, Blue, ""))+"SQL  "+Reset, fileAndLine, showStr)
}

func Info(args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(fmt.Sprintf("%s", util.If(isColor, Magenta, ""))+"Info "+Reset, fmt.Sprintf("%s:%d", file, line), args...)
}

func Println(args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon("Info ", fmt.Sprintf("%s:%d", file, line), args...)
}

func Warn(args ...interface{}) {
	if logLevel < WarnLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(fmt.Sprintf("%s", util.If(isColor, Yellow, ""))+"Warn "+Reset, fmt.Sprintf("%s:%d", file, line), args...)
}

func Error(args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	logCommon(fmt.Sprintf("%s", util.If(isColor, Red, ""))+"Error"+Reset, fmt.Sprintf("%s:%d", file, line), args...)
}

func Fatal(args ...interface{}) {
	if logLevel < FatalLevel {
		return
	}
	log.Fatal(args)
}

func Panic(args ...interface{}) {
	if logLevel < PanicLevel {
		return
	}
	log.Panic(args)
}

func InfoByteListHex(byteList []byte) {
	str := "["
	for i, b := range byteList {
		str += fmt.Sprintf(" 0x%02X", b)
		if i == len(byteList)-1 {
			str += " "
		}
	}
	str += "]"
	Info(str)
}

func WithFields(msg string, fields Fields, level Level) {
	_, file, line, _ := runtime.Caller(1)
	file = strings.Replace(file, logMainPrefix, "", 1)

	var fieldsJson string
	if fields != nil {
		fieldsJsonByte, err := json.Marshal(fields)
		if err != nil {
			fieldsJson = fmt.Sprintf("%v", fields)
		} else {
			fieldsJson = string(fieldsJsonByte)
		}
	}
	levelStr := ""
	switch level {
	case TraceLevel:
		levelStr = "Trace"
	case DebugLevel:
		levelStr = "Debug"
	case InfoLevel:
		levelStr = "Info "
	case WarnLevel:
		levelStr = "Warn "
	case ErrorLevel:
		levelStr = "Error"
	case FatalLevel:
		levelStr = "Fatal"
	case PanicLevel:
		levelStr = "Panic"
	}
	logCommon(levelStr, fmt.Sprintf("%s:%d", file, line), msg, fieldsJson)
}
