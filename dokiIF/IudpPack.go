package dokiIF

import "net"

type IUdpDataPack interface {
	//封包的方法
	Pack(uint32, []byte) ([]byte, error)
	//拆包的方法
	UnPack([]byte, *net.UDPAddr) (uint32, []byte, error)
}
