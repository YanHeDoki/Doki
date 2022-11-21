package dokiIF

import "net"

//IRequest
//实际上是吧客户端请求的链接信息，和请求的数据 包装到了一个Request中

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection
	//upd连接操作
	GetUdpConn() *net.UDPConn
	//得到请求的消息数据
	GetData() []byte
	//得到消息id
	GetMsgId() uint32

	BindRouter(IRouter)
	//执行下一个函数
	Next()
	//终结路由函数的执行
	Abort()
	//是否终结了函数
	IsAbort() bool
}
