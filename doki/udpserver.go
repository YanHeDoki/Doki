package doki

import (
	"fmt"
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
)

type UdpServer struct {
	//服务器ip
	IP string
	//服务器监听的端口
	Port int
	//服务器连接
	UdpConn *net.UDPConn
	//通知服务器退出
	exitChan chan struct{}
	//当前的对象添加一个router server注册的链接对应的业务
	//当前Server的消息管理模块，用来绑定MsgId和对应的router
	MsgHandler dokiIF.IUdpMsgHandle
	//拆封包工具
	packet dokiIF.IUdpDataPack
}

func NewUdpServer(IP string, Prot int, pack dokiIF.IUdpDataPack) dokiIF.IUdpServer {
	return &UdpServer{
		IP:         IP,
		Port:       Prot,
		exitChan:   nil,
		MsgHandler: NewUdpMsgHandle(),
		packet:     pack,
	}
}

func (u *UdpServer) Start() {
	//日志，以后应该用日志来处理
	BaseLog.DefaultLog.DokiLog("info", fmt.Sprintf("[START] Server: %s Listener at IP: %s  Port: %d starting", u.IP, u.Port))
	u.exitChan = make(chan struct{})

	go func() {
		//获取一个Tcp的Addr地址
		resolveUdpPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", u.IP, u.Port))
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Start ServerErr err:%s", err))
			panic(err)
		}
		//监听服务器的地址
		listen, err := net.ListenUDP("udp", resolveUdpPAddr)
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("ListenIPErr err:%s", err))
			panic(err)
		}
		//绑定自己的udp连接
		u.UdpConn = listen
		BaseLog.DefaultLog.DokiLog("info", "Start  Udp Server success Listening...")
		//另起协程去监听
		go func() {
			//阻塞的等待客户端的连接 处理客户端的链接业务（读写）
			for {
				ReadUdp := make([]byte, 4096)
				n, addr, err := listen.ReadFromUDP(ReadUdp)
				if err != nil {
					BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("ReadFromUDP err:%s", err))
					continue
				}
				id, data, err := u.GetPacket().UnPack(ReadUdp[:n], addr)
				if err != nil {
					BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("UnPack err:%s", err))
					continue
				}
				go u.MsgHandler.DoMsgHandler(&UdpRequest{
					id:      id,
					data:    data,
					udpConn: listen,
				})
			}
		}()
		//阻塞用来通知退出
		select {
		case <-u.exitChan:
			err := listen.Close()
			if err != nil {
				BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Listener close err:%s ", err))
			}
		}

	}()
}

func (u *UdpServer) Stop() {
	BaseLog.DefaultLog.DokiLog("info", "[STOP]  Udpserver")
	//断开服务器，将一些服务器的资源链接释放
	u.exitChan <- struct{}{}
	close(u.exitChan)
}

func (u *UdpServer) Server() {
	u.Start()
	select {}
}

func (u *UdpServer) AddRouter(msgid uint32, router ...dokiIF.RouterHandler) {
	u.MsgHandler.AddRouter(msgid, router...)
}

func (u *UdpServer) GetPacket() dokiIF.IUdpDataPack {
	return u.packet
}

func (u *UdpServer) GetUdpConn() *net.UDPConn {
	return u.UdpConn
}
