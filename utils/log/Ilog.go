package BaseLog

type Ilog interface {
	DokiLog(LogLv, format string, arg ...interface{})
}
