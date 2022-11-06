package iface

type Inotify interface {
	//通知某个id的方法
	NotifyToConnByID(Id uint64, msg IMessage)
	//通知所有人
	NotifyAll(msg IMessage)
}
