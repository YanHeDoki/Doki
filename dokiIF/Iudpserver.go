package dokiIF

import "net"

//定义一个服务器接口

type IUdpServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Server()

	//udp服务连接
	GetUdpConn() *net.UDPConn
	//新路由功能：给当前的服务器注册一个路由方法。供客户端的链接处理使用
	AddRouter(msgid uint32, router ...RouterHandler)
	//支持自定义的封包
	GetPacket() IDataPack
}
