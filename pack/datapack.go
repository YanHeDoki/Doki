package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/conf"
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"io"
)

//DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

var defaultHeaderLen uint32 = 8

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节)+ID uint32（4字节）
	return defaultHeaderLen
}

//Pack 封包方法(压缩数据)
func (d *DataPack) Pack(msg dokiIF.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将data id 封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data长度封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//将data 消息本体封装进去
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//UnPack 拆包方法 （将包的head信息读取）之后再根据head里的len信息读取信息
func (d *DataPack) UnPack(conn io.Reader) (dokiIF.IMessage, error) {

	headBuff := make([]byte, d.GetHeadLen())
	_, err := io.ReadFull(conn, headBuff)
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("read in packhead err:%s", err.Error()))
		return nil, err
	}

	//创建一个从二进制读取数据的ioReader
	dataBuff := bytes.NewReader(headBuff)

	//只解压head信息的到len和id
	msg := &Message{}
	//读取id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//读取datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//判断datalen是否过长超出允许长度
	if conf.GlobalConfObject.MaxPacketSize > 0 && msg.DataLen > conf.GlobalConfObject.MaxPacketSize {
		return nil, errors.New("too large msg data recv!")
	}

	//根据datalen的参数再去读取一次
	var data []byte
	if msg.GetDataLen() > 0 {
		data = make([]byte, msg.GetDataLen())
		if _, err := io.ReadFull(conn, data); err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("read data err:%s", err.Error()))
			return nil, err
		}
	}
	msg.SetData(data)
	return msg, nil
}
