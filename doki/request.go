package doki

import "github.com/YanHeDoki/Doki/iface"

type Request struct {
	//已经和客户端建立好的链接
	conn iface.IConnection
	//客户端请求的数据
	msg iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
