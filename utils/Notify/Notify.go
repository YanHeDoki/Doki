package Notify

import (
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"sync"
)

type notify struct {
	cimap ConnIDMap
	sync.RWMutex
}

func NewNotify() *notify {
	return &notify{
		cimap: make(map[uint64]dokiIF.IConnection, 1000),
	}
}

func (n *notify) HasIdConn(Id uint64) bool {
	n.RLock()
	defer n.RUnlock()
	_, ok := n.cimap[Id]
	return ok
}

func (n *notify) SetNotifyID(Id uint64, conn dokiIF.IConnection) {
	n.Lock()
	defer n.Unlock()
	n.cimap[Id] = conn
}

func (n *notify) GetNotifyByID(Id uint64) (dokiIF.IConnection, error) {
	n.RLock()
	defer n.RUnlock()
	Conn, ok := n.cimap[Id]
	if !ok {
		return nil, errors.New(" Not Find UserId")
	}
	return Conn, nil
}

func (n *notify) DelNotifyByID(Id uint64) {
	n.RLock()
	defer n.RUnlock()
	delete(n.cimap, Id)
}

func (n *notify) NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, err := n.GetNotifyByID(Id)
	if err != nil {
		return err
	}
	err = Conn.SendMsg(MsgId, data)
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Notify to %d err:%s", Id, err))
		return err
	}
	return nil
}

func (n *notify) NotifyAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendMsg(MsgId, data)
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Notify to %d err:%s", Id, err))
			return err
		}
	}
	return nil
}

func (n *notify) notifyAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendMsg(MsgId, data)
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Notify to %d err:%s", Id, err))
			return err
		}
	}
	return nil
}

//极端情况 同时加入和发送的人很多需要尽快释放map的情况， 目前问题很多不采用
func (n *notify) notifyAll2(MsgId uint32, data []byte) error {
	conns := make([]dokiIF.IConnection, 0, len(n.cimap))
	n.RLock()
	for _, v := range n.cimap {
		conns = append(conns, v)
	}
	n.RUnlock()

	for i := 0; i < len(conns); i++ {
		conns[i].SendMsg(MsgId, data)
	}
	return nil
}

func (n *notify) NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, err := n.GetNotifyByID(Id)
	if err != nil {
		return err
	}
	err = Conn.SendBuffMsg(MsgId, data)
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Notify to %d err:%s", Id, err))
		return err
	}
	return nil
}

func (n *notify) NotifyBuffAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendBuffMsg(MsgId, data)
		if err != nil {
			BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Notify to %d err:%s", Id, err))
			return err
		}
	}
	return nil
}
