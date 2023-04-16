package doki

import (
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/conf"
	"github.com/YanHeDoki/Doki/constants"
	"github.com/YanHeDoki/Doki/dokiIF"
	"github.com/YanHeDoki/Doki/pack"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
)

var logo = `
 _____        _     _     _       _     _ 
(____ \      | |   (_)   | |     | |   (_)
 _   \ \ ___ | |  _ _  _ | | ___ | |  _ _ 
| |   | / _ \| | / ) |/ || |/ _ \| | / ) |
| |__/ / |_| | |< (| ( (_| | |_| | |< (| |
|_____/ \___/|_| \_)_|\____|\___/|_| \_)_|
`

type Server struct {
	//服务器名称
	Name string
	//IP版本 IPv4 or other
	IPVersion string
	//服务器ip
	IP string
	//服务器监听的端口
	Port int
	//通知服务器退出
	exitChan chan struct{}
	//当前的对象添加一个router server注册的链接对应的业务
	//当前Server的消息管理模块，用来绑定MsgId和对应的router
	MsgHandler dokiIF.IMsgHandle
	//该server的连接管理器
	ConnMgr dokiIF.IConnManager

	//新增两个hook函数原型
	//该Server的连接创建时Hook函数
	OnConnStart func(conn dokiIF.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn dokiIF.IConnection)

	//拆封包工具
	packet dokiIF.IDataPack
}

// DefaultServer 初始化默认server服务器方法
func DefaultServer() dokiIF.IServer {
	//读取配置
	conf.ConfigInit()
	//打印logo
	printLogo()
	BaseLog.DefaultLog = BaseLog.NewLog("info")
	return &Server{ //报错不能返回这个类型
		Name:       conf.GlobalConfObject.Name,
		IPVersion:  "tcp",
		IP:         conf.GlobalConfObject.Host,
		Port:       conf.GlobalConfObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
		packet:     pack.Factory().NewPack(conf.GlobalConfObject.PacketName),
	}
}

func NewServer(config *conf.Config) dokiIF.IServer {
	//注入用户配置
	conf.UserConfInit(config)
	//打印logo
	printLogo()
	s := &Server{
		Name:       config.Name,
		IPVersion:  config.TcpVersion,
		IP:         config.Host,
		Port:       config.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
	}
	if config.UserPack == nil {
		s.packet = pack.Factory().NewPack(constants.StdDataPack)
	} else {
		s.packet = config.UserPack
	}
	if config.Log == nil {
		BaseLog.DefaultLog = BaseLog.NewLog(config.LogLevel)
	} else {
		BaseLog.DefaultLog = config.Log
	}

	return s
}

func (s *Server) Start() {
	//日志，以后应该用日志来处理
	BaseLog.DefaultLog.DokiLog("info", fmt.Sprintf("[START] Server: %s Listener at IP: %s  Port: %d starting", s.Name, s.IP, s.Port))
	s.exitChan = make(chan struct{})
	//开启工作线程池
	s.MsgHandler.StartWorkerPool()
	//由server方法来阻塞所以异步处理

	go func() {
		//获取一个Tcp的Addr地址
		resolveIPAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Start ServerErr err:%s", err))
			panic(err)
		}
		//监听服务器的地址
		listen, err := net.ListenTCP(s.IPVersion, resolveIPAddr)
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("ListenIPErr err:%s", err))
			panic(err)
		}

		BaseLog.DefaultLog.DokiLog("info", fmt.Sprintf("Start  %s success Listening...", s.Name))
		var cid uint32
		cid = 0

		//另起协程去监听
		go func() {
			//阻塞的等待客户端的连接 处理客户端的链接业务（读写）
			for {
				conn, err := listen.AcceptTCP()
				if err != nil {
					//Go 1.16+ 判断是否是net.ErrClosed 既监听链接已经关闭
					if errors.Is(err, net.ErrClosed) {
						BaseLog.DefaultLog.DokiLog("error", "Listener closed")
						return
					}
					BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("AcceptTCP err:%s", err))
					continue
				}
				//设置最大连接数量的判断，如果超过最大连接数就断开
				if s.ConnMgr.Len() >= conf.GlobalConfObject.MaxConn {
					//todo 给客户端一个错误信息
					bytes, err := s.packet.Pack(pack.NewMsgPackage(0, []byte("Server Conn is Max....")))
					if err != nil {
						BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Server MaxConn send msg err:%s", err))
					}
					conn.Write(bytes)
					BaseLog.DefaultLog.DokiLog("warning", fmt.Sprintf("Too Many Connections MaxConn=%d", conf.GlobalConfObject.MaxConn))
					conn.Close()
					continue
				}
				//使用新的connection模块
				newConnection := NewConnection(s, conn, cid)

				cid++
				go newConnection.Start()
			}
		}()
		//阻塞用来通知退出
		select {
		case <-s.exitChan:
			err := listen.Close()
			if err != nil {
				BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Listener close err:%s ", err))
			}
		}

	}()

}

func (s *Server) Stop() {
	BaseLog.DefaultLog.DokiLog("info", fmt.Sprintf("[STOP]  server , name:%s", s.Name))
	//断开服务器，将一些服务器的资源链接释放
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

func (s *Server) Server() {

	//启动服务器
	s.Start()

	//TODO 留空位可以给以后操作空间
	//阻塞 否则主Go退出， listenner的go将会退出
	select {}
}

// 新方法
func (s *Server) AddRouter(msgId uint32, router ...dokiIF.RouterHandler) dokiIF.IRouter {
	return s.MsgHandler.AddRouter(msgId, router...)
}

// 路由组管理
func (s *Server) Group(start, end uint32, Handlers ...dokiIF.RouterHandler) dokiIF.IGroupRouter {
	return s.MsgHandler.Group(start, end, Handlers...)
}
func (s *Server) Use(Handlers ...dokiIF.RouterHandler) dokiIF.IRouter {
	return s.MsgHandler.Use(Handlers...)
}

func (s *Server) GetConnMgr() dokiIF.IConnManager {
	return s.ConnMgr
}

func (s *Server) GetMsgHandler() dokiIF.IMsgHandle {
	return s.MsgHandler
}

func (s *Server) SetOnConnStart(hookFunc func(dokiIF.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(dokiIF.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn dokiIF.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn dokiIF.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) GetPacket() dokiIF.IDataPack {
	return s.packet
}

func printLogo() {
	fmt.Println(logo)
	fmt.Printf("[Doki] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		conf.GlobalConfObject.Version,
		conf.GlobalConfObject.MaxConn,
		conf.GlobalConfObject.MaxPacketSize)
}
