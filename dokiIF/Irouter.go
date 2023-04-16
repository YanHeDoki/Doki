package dokiIF

// 路由接口
// 路由里的接口都是IRequest
type RouterHandler func(request IRequest)

type IRouter interface {
	//添加全局组件
	Use(Handlers ...RouterHandler)
	//添加路由
	AddHandler(msgId uint32, handlers ...RouterHandler)

	//路由组管理
	Group(start, end uint32, Handlers ...RouterHandler) IGroupRouter

	GetHandlers(MsgId uint32) ([]RouterHandler, bool)
	//执行下一个函数
	//Next(request IRequest)
	////终结路由函数的执行
	//Abort()
	////是否终结了函数
	//IsAbort() bool
}

type IGroupRouter interface {
	//添加全局组件
	Use(Handlers ...RouterHandler)
	//添加组路由组件
	AddHandler(MsgId uint32, Handlers ...RouterHandler)
}
