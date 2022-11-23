package BaseLog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

// 定义日志等级，使用Uint16定义一个类型来进行比较的支持
type Level uint16

//定义日志等级,常量定义
const (
	UNKONW Level = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

var DefaultLog Ilog

//日志结构体
type Logger struct {
	LogLevel Level
}

func (l *Logger) DokiLog(LogLv, LogStr string) {
	switch LogLv {
	case "debug":
		l.Debug(LogStr)
	case "info":
		l.Info(LogStr)
	case "warning":
		l.Warning(LogStr)
	case "error":
		l.Error(LogStr)
	}
}

//把传入的等级转换成可以比较的等级
func paserLoglevel(s string) (Level, error) {
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
	default:
		err := fmt.Errorf("无效级别")
		return UNKONW, err
	}
}

//反向转换日志等级为string
func paserLogString(logLevel Level) (Level string) {

	switch logLevel {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		fmt.Println("无效级别")
		return "UNKONW"
	}
}

//构造函数
//调用转换函数,把传入的string类型转成自定义的Level等级进行比较
func NewLog(lvstr string) *Logger {

	level, err := paserLoglevel(lvstr)
	if err != nil {
		fmt.Println("传入参数出错:", err)
		panic(err)
	}
	return &Logger{
		LogLevel: level,
	}
}

//比较方法，用于比较需要打印的和传入构造的等级
func (l *Logger) enable(level Level) bool {
	return l.LogLevel <= level

}

//日志打印方法
func (l *Logger) log(lv Level, logStr string) {
	if l.enable(lv) {
		t := time.Now()
		_, filename, lineNo := getInfo(4)
		fmt.Printf("[%s] [%s][%s:%d] %s \n", t.Format("2006-01-02 15:04:05"), paserLogString(lv), filename, lineNo, logStr)
	}
}

//日志等级调用打印方法
func (l *Logger) Debug(logStr string) {
	if l.enable(DEBUG) {
		l.log(DEBUG, logStr)
	}
}

func (l *Logger) Info(logStr string) {

	l.log(INFO, logStr)

}

func (l *Logger) Warning(logStr string) {

	l.log(WARNING, logStr)

}

func (l *Logger) Error(logStr string) {

	l.log(ERROR, logStr)

}

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
