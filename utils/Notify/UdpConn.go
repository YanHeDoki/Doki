package Notify

import (
	"errors"
	"github.com/YanHeDoki/Doki/dokiIF"
	"github.com/YanHeDoki/Doki/pack"
	"net"
	"sync"
)

type UdpConnMap map[uint64]*net.UDPAddr

type UdpNotify struct {
	UCM       UdpConnMap
	udpServer dokiIF.IUdpServer
	sync.RWMutex
}

func NewUdpNotify(udpserver dokiIF.IUdpServer) *UdpNotify {
	return &UdpNotify{
		UCM:       make(UdpConnMap, 100),
		udpServer: udpserver,
	}
}

func (un *UdpNotify) HasId(Id uint64) bool {
	un.RLock()
	defer un.RUnlock()
	_, ok := un.UCM[Id]
	return ok
}

func (un *UdpNotify) SetAddr(Id uint64, addr *net.UDPAddr) {
	un.Lock()
	defer un.Unlock()
	un.UCM[Id] = addr
}

func (un *UdpNotify) GetAddrById(Id uint64) (*net.UDPAddr, error) {
	un.RLock()
	defer un.RUnlock()
	addr, ok := un.UCM[Id]
	if !ok {
		return nil, errors.New("Not Find Id Addr")
	}
	return addr, nil
}

func (un *UdpNotify) DelAddrById(Id uint64) {
	un.Lock()
	defer un.Unlock()
	delete(un.UCM, Id)
}

func (un *UdpNotify) SendUdp(Id uint64, data []byte) error {
	un.RLock()
	defer un.RUnlock()
	_, err := un.udpServer.GetUdpConn().WriteToUDP(data, un.UCM[Id])
	return err
}

func (un *UdpNotify) SendUdpTo(Id uint64, MsgId uint32, data []byte) error {
	un.RLock()
	addr := un.UCM[Id]
	un.RUnlock()
	resp, err := un.udpServer.GetPacket().Pack(pack.NewMsgPackage(MsgId, data))
	if err != nil {
		return err
	}
	_, err = un.udpServer.GetUdpConn().WriteToUDP(resp, addr)
	if err != nil {
		return err
	}
	return nil
}

func (un *UdpNotify) Broadcast(Ids []uint64, MsgId uint32, data []byte) error {

	resp, err := un.udpServer.GetPacket().Pack(pack.NewMsgPackage(MsgId, data))
	if err != nil {
		return err
	}
	for _, v := range Ids {
		err := un.SendUdp(v, resp)
		if err != nil {
			return err
		}
	}
	return nil
}
