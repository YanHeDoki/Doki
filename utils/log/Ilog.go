package BaseLog

type Ilog interface {
	DokiLog(LogLv, LogStr string)
}
