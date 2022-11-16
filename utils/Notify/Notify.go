package Notify

import (
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/dokiIF"
	"sync"
)

type notify struct {
	cimap ConnIDMap
	look  sync.RWMutex
}

func NewNotify() *notify {
	return &notify{
		cimap: make(map[uint64]dokiIF.IConnection, 100),
	}
}

func (n *notify) SetNotifyID(Id uint64, conn dokiIF.IConnection) {
	n.look.Lock()
	defer n.look.Unlock()
	n.cimap[Id] = conn
}

func (n *notify) GetNotifyByID(Id uint64) (dokiIF.IConnection, error) {
	n.look.RLock()
	defer n.look.RUnlock()
	Conn, ok := n.cimap[Id]
	if !ok {
		return nil, errors.New(" Not Find UserId")
	}
	return Conn, nil
}

func (n *notify) DelNotifyByID(Id uint64) {
	n.look.RLock()
	defer n.look.RUnlock()
	delete(n.cimap, Id)
}

func (n *notify) NotifyToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, ok := n.cimap[Id]
	if !ok {
		return errors.New(" Not Find UserId")
	}
	err := Conn.SendMsg(MsgId, data)
	if err != nil {
		fmt.Println("Notify to", Id, "err:", err)
		return err
	}
	return nil
}

func (n *notify) NotifyAll(MsgId uint32, data []byte) error {
	for id, v := range n.cimap {
		err := v.SendMsg(MsgId, data)
		if err != nil {
			fmt.Println("Notify to", id, "err:", err)
			return err
		}
	}
	return nil
}

func (n *notify) NotifyBuffToConnByID(Id uint64, MsgId uint32, data []byte) error {
	Conn, ok := n.cimap[Id]
	if !ok {
		return errors.New(" Not Find UserId")
	}
	err := Conn.SendBuffMsg(MsgId, data)
	if err != nil {
		fmt.Println("Notify to", Id, "err:", err)
		return err
	}
	return nil
}

func (n *notify) NotifyBuffAll(MsgId uint32, data []byte) error {
	for id, v := range n.cimap {
		err := v.SendBuffMsg(MsgId, data)
		if err != nil {
			fmt.Println("Notify to", id, "err:", err)
			return err
		}
	}
	return nil
}
