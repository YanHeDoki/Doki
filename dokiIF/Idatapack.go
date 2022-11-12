package dokiIF

import (
	"io"
)

type IDataPack interface {
	//获取包的头的长度方法
	GetHeadLen() uint32
	//封包的方法
	Pack(msg IMessage) ([]byte, error)
	//拆包的方法
	UnPack(reader io.Reader) (IMessage, error)
}
