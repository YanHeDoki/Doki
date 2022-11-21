package conf

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
)

type Config struct {

	//Server
	TcpServer  dokiIF.IServer //当前全局的Server对象
	Host       string         //当前服务器主机监听的IP
	TcpPort    int            //当前服务器监听的端口
	Name       string         //当前服务器的名称
	TcpVersion string         //tcp版本

	//服务器可选配置
	Version          string           //版本
	MaxConn          int              //最大连接数量
	MaxPacketSize    uint32           //当前框架数据包的最大尺寸
	WorkerPoolSize   uint32           //业务工作Worker池的数量
	MaxWorkerTaskLen uint32           //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32           //SendBuffMsg发送消息的缓冲最大长度
	UserPack         dokiIF.IDataPack //用户自定义封解包

	Log BaseLog.Ilog
}
