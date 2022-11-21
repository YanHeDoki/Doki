package pack

import (
	"bytes"
	"encoding/binary"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
)

type UdpPack struct {
}

//Pack 封包方法(压缩数据)
func (d *UdpPack) Pack(id uint32, data []byte) ([]byte, error) {
	BaseLog.DefaultLog.DokiLog("warning", "please implement UdpPack Not Use This Func")
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将data id 封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, id); err != nil {
		return nil, err
	}
	//将data 消息本体封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//UnPack 拆包方法 （将包的head信息读取）之后再根据head里的len信息读取信息
func (d *UdpPack) UnPack(data []byte, addr *net.UDPAddr) (uint32, []byte, error) {

	BaseLog.DefaultLog.DokiLog("warning", "please implement UdpPack Not Use This Func")
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	err := binary.Read(dataBuff, binary.LittleEndian, data)
	if err != nil {
		return 0, nil, err
	}
	return 0, dataBuff.Bytes(), nil
}
