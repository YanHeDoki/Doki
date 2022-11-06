package iface

//路由接口
//路由里的接口都是IRequest

//尝试修改中
type RouterHandler func(router IRouter, request IRequest)
type IRouter interface {
	//执行下一个函数
	Next(request IRequest)
	//终结路由函数的执行
	Abort()
	//是否终结了函数
	IsAbort() bool
}
