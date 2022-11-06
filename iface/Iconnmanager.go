package iface

//连接管理模块

type IConnManager interface {

	//添加用户链接
	Add(connection IConnection)
	//删除链接
	Remove(connection IConnection)
	//根据connid获取一个链接
	Get(connId uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清楚并终止所有链接
	ClearConn()
}
