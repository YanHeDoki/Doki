package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
)

type Request struct {
	//已经和客户端建立好的链接
	conn dokiIF.IConnection
	//客户端请求的数据
	msg dokiIF.IMessage
	//组合路由接口
	Router dokiIF.IRouter
}

func (r *Request) GetUdpConn() *net.UDPConn {
	BaseLog.DefaultLog.DokiLog("warning", "Udp TcpRequest Is Not have UDPConn")
	return nil
}

func (r *Request) GetConnection() dokiIF.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) BindRouter(router dokiIF.IRouter) {
	r.Router = router
}

func (r *Request) Next() {
	r.Router.Next(r)
}

func (r *Request) Abort() {
	r.Router.Abort()
}

func (r *Request) IsAbort() bool {
	return r.Router.IsAbort()
}

type UdpRequest struct {
	id      uint32
	data    []byte
	udpConn *net.UDPConn
	Router  dokiIF.IRouter
}

func (u *UdpRequest) GetUdpConn() *net.UDPConn {
	return u.udpConn
}

func (u *UdpRequest) GetConnection() dokiIF.IConnection {
	BaseLog.DefaultLog.DokiLog("warning", "Udp UdpRequest Is Not have Connection")
	return nil
}

func (u *UdpRequest) BindRouter(router dokiIF.IRouter) {
	u.Router = router
}

func (u *UdpRequest) Next() {
	u.Router.Next(u)
}

func (u *UdpRequest) Abort() {
	u.Router.Abort()
}

func (u *UdpRequest) IsAbort() bool {
	return u.Router.IsAbort()
}

func (u *UdpRequest) GetData() []byte {
	return u.data
}

func (u *UdpRequest) GetMsgId() uint32 {
	return u.id
}
