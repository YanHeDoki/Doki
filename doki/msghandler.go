package doki

import (
	"fmt"
	"github.com/YanHeDoki/Doki/conf"
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"sync"
)

type MsgHandle struct {
	Router         *Router                //路由模块
	WorkerPoolSize uint32                 //业务工作Worker池的数量
	TaskQueue      []chan dokiIF.IRequest //Worker负责取任务的消息队列
	sync.RWMutex
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Router:         NewRouter(),
		WorkerPoolSize: conf.GlobalConfObject.WorkerPoolSize, //注意一个消息队列对应一个worker池子
	}
}

// 尝试修改msghandler
func (m *MsgHandle) DoMsgHandler(request dokiIF.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("doMsgHandler panic %v :", err))
		}
	}()
	handlers, ok := m.Router.GetHandlers(request.GetMsgId())
	if !ok {
		BaseLog.DefaultLog.DokiLog("warning", fmt.Sprintf("api msgID = %d is not FOUND!", request.GetMsgId()))
		return
	}
	request.BindRouter(handlers)
	request.Next()
}

func (m *MsgHandle) AddRouter(msgId uint32, handler ...dokiIF.RouterHandler) dokiIF.IRouter {
	m.Router.AddHandler(msgId, handler...)
	return m.Router
}

func (m *MsgHandle) Group(start, end uint32, Handlers ...dokiIF.RouterHandler) dokiIF.IGroupRouter {
	return NewGroup(start, end, m.Router, Handlers...)
}
func (m *MsgHandle) Use(Handlers ...dokiIF.RouterHandler) dokiIF.IRouter {
	m.Router.Use(Handlers...)
	return m.Router
}

func (m *MsgHandle) StartWorkerPool() {

	if m.TaskQueue == nil {
		//优化内存占用是tcp会主动调用这个方法的时候才分配内存udp服务不开启队列消息所以不需要开辟内存
		//只有开启的时候才分配内存
		m.TaskQueue = make([]chan dokiIF.IRequest, conf.GlobalConfObject.WorkerPoolSize)
	}

	//根据配置的workerpool的size来分别开启worker 每个worker用一个go承载
	for i := uint32(0); i < m.WorkerPoolSize; i++ {
		//一个worker被启动
		//1.当前的worker对应的channel消息队列 开辟对应的空间 0号worker对应0号channel
		//用MaxWorkerTaskLen限制一个管道最多接受多少条消息
		m.TaskQueue[i] = make(chan dokiIF.IRequest, conf.GlobalConfObject.MaxWorkerTaskLen)
		go m.startOneWorker(i)
	}

}

func (m *MsgHandle) startOneWorker(workerId uint32) {

	//不断的阻塞去等代消息
	for i := 0; i < conf.GlobalConfObject.DoMsgHandlerNum; i++ {
		go func() {
			for {
				select {
				//根据id去结构体中取到对应的消息队列来消费，如果管道中有消息的话
				case req := <-m.TaskQueue[workerId]:
					m.DoMsgHandler(req)
				}
			}
		}()
	}

}

func (m *MsgHandle) SendMsgToTaskQueue(request dokiIF.IRequest) {
	//将消息平均的分配给woroker
	//根据客户端建立的连接id来判断
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	//将消息发送给消息队列
	m.TaskQueue[workerId] <- request
}
