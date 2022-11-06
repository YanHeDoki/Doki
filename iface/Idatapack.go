package iface

import "io"

type IDataPack interface {
	//获取包的头的长度方法
	GetHeadLen() uint32
	//封包的方法
	Pack(msg IMessage) ([]byte, error)
	//拆包的方法
	UnPack(reader io.Reader) (IMessage, error)
}

const (
	//标准封包和拆包方式
	StdDataPack string = "std_pack"

	//自定义封包方式在此添加
)

const (
	//默认标准报文协议格式
	Message string = "message"
)
