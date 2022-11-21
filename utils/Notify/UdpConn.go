package Notify

import (
	"errors"
	"net"
	"sync"
)

type UdpConnMap map[uint64]*net.UDPAddr

type UdpNotify struct {
	UCM     UdpConnMap
	UdpConn *net.UDPConn
	sync.RWMutex
}

func (un *UdpNotify) HasId(Id uint64) bool {
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

func (un *UdpNotify) SendUdpTo(Id uint64, data []byte) error {
	un.Lock()
	defer un.Unlock()
	_, err := un.UdpConn.WriteToUDP(data, un.UCM[Id])
	if err != nil {
		return err
	}
	return nil
}

func (un *UdpNotify) Broadcast(Ids []uint64, data []byte) error {
	un.Lock()
	defer un.Unlock()
	for _, v := range Ids {
		_, err := un.UdpConn.WriteToUDP(data, un.UCM[v])
		if err != nil {
			return err
		}
	}
	return nil
}
