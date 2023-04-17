package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
)

type Request struct {
	//已经和客户端建立好的链接
	conn dokiIF.IConnection
	//Udp连接
	udpConn *net.UDPConn
	//客户端请求的数据
	msg dokiIF.IMessage
	//组合路由操作
	handlers []dokiIF.RouterHandler
	index    int8
}

func (r *Request) GetUdpConn() *net.UDPConn {
	if r.udpConn == nil {
		BaseLog.DefaultLog.DokiLog("warning", "TcpRequest Is Not have UDPConn")
		return nil
	}
	return r.udpConn
}

func (r *Request) GetConnection() dokiIF.IConnection {
	if r.conn == nil {
		BaseLog.DefaultLog.DokiLog("warning", "UdpRequest Is Not have TcpConn")
		return nil
	}
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) BindRouter(handlers []dokiIF.RouterHandler) {
	r.handlers = handlers
	r.index = -1
}

func (r *Request) Next() {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r)
		r.index++
	}
}

func (r *Request) Abort() {
	r.index = int8(len(r.handlers))
}

func (r *Request) IsAbort() bool {
	return r.index >= int8(len(r.handlers))
}
