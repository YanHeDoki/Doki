package log

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type LogLevel int16

//定义日志等级,常量定义
const (
	UNKNOWN LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

func StartFile(lv, fp, fn, model string, logH int, maxsize int64) *filelog {
	return NewFilelog(lv, fp, fn, model, logH, maxsize)
}

//把传入的等级转换成可以比较的等级
func paserLoglevel(s string) (LogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil

	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := fmt.Errorf("无效级别")
		return UNKNOWN, err
	}
}

//反向转换日志等级为string
func paserLogString(logLevel LogLevel) (Level string) {

	switch logLevel {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		fmt.Println("无效级别")
		return "UNKNOWN"
	}

}

//获取文件信息
func getInfo(ship int) (funcname, filename string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(ship)
	if !ok {
		fmt.Println("runtime.caller() err")
		return
	}
	funcname = runtime.FuncForPC(pc).Name()
	filename = path.Base(file)
	funcname = strings.Split(funcname, ".")[1]
	return funcname, filename, lineNo

}
