package dokiIF

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	DoMsgHandler(request IRequest) //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, handler ...RouterHandler) IRouter
	Group(start, end uint32, Handlers ...RouterHandler) IGroupRouter
	Use(Handlers ...RouterHandler) IRouter
	StartWorkerPool()                //启动工作线程池
	SendMsgToTaskQueue(req IRequest) //将消息发送到消息队列去处理
}
