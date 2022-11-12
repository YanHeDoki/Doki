package dokiIF

//定义一个服务器接口

type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Server()

	//新路由功能：给当前的服务器注册一个路由方法。供客户端的链接处理使用
	AddRouter(msgid uint32, router ...RouterHandler)

	//返回连接资源管理器
	GetConnMgr() IConnManager
	//GetMsgHandler 获取MsgHandler管理器避免再复制一次
	GetMsgHandler() IMsgHandle
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
	//支持自定义的封包
	GetPacket() IDataPack
}
