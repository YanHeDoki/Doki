package pack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"io"
	"net"
)

type UdpPack struct {
}

func (d *UdpPack) GetHeadLen() uint32 {
	return 0
}

func (d *UdpPack) Pack(msg dokiIF.IMessage) ([]byte, error) {
	BaseLog.DefaultLog.DokiLog("warning", "please implement UdpPack Not Use This Func")
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将data id 封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data 消息本体封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (d *UdpPack) UnPack(conn io.Reader) (dokiIF.IMessage, error) {
	BaseLog.DefaultLog.DokiLog("warning", "please implement UdpPack Not Use This Func")

	ReadUdp := make([]byte, 4096)
	udpconn := conn.(*net.UDPConn)
	n, _, err := udpconn.ReadFromUDP(ReadUdp)
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("ReadFromUDP err:%s", err))
		return nil, err
	}
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	err = binary.Read(dataBuff, binary.LittleEndian, ReadUdp[:n])
	if err != nil {
		return nil, err
	}
	return NewMsgPackage(0, dataBuff.Bytes()), nil
}
