package doki

type Inotify interface {
	//通知某个id的方法
	NotifyToConnByID(Id uint64, msg IMessage) error
	//通知所有人
	NotifyAll(msg IMessage) error

	//通知某个id的方法
	NotifyBuffToConnByID(Id uint64, msg IMessage) error
	//通知所有人
	NotifyBuffAll(msg IMessage) error
}
