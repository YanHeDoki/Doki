package iface

import (
	"context"
	"net"
)

type IConnection interface {

	//开始链接方法
	Start()
	//停止链接方法
	Stop()
	//获取当前链接绑定的socket
	GetTcpConnection() *net.TCPConn
	//返回ctx，用于用户自定义的go程获取连接退出状态
	GetContext() context.Context
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的Tcp状态 ip port
	RemoteAddr() net.Addr
	//直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendMsg(msgID uint32, data []byte) error
	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error //添加带缓冲发送消息接口
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}
