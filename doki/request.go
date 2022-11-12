package doki

import "github.com/YanHeDoki/Doki/dokiIF"

type Request struct {
	//已经和客户端建立好的链接
	conn dokiIF.IConnection
	//客户端请求的数据
	msg dokiIF.IMessage
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
