package Notify

import (
	"fmt"
	"github.com/YanHeDoki/Doki/doki"
	"sync"
)

type notify struct {
	cimap ConnIDMap
	look  sync.RWMutex
}

func NewNotify() *notify {
	return &notify{
		cimap: make(map[uint64]doki.IConnection, 100),
	}
}

func (n *notify) SetNotifyID(Id uint64, conn doki.IConnection) {
	n.look.Lock()
	defer n.look.RLock()
	n.cimap[Id] = conn
}
func (n *notify) GetNotifyByID(Id uint64) doki.IConnection {
	n.look.RLock()
	defer n.look.RLock()
	return n.cimap[Id]
}

func (n *notify) NotifyToConnByID(Id uint64, msg doki.IMessage) error {
	err := n.cimap[Id].SendMsg(msg.GetMsgId(), msg.GetData())
	if err != nil {
		fmt.Println("Notify to", Id, "err:", err)
		return err
	}
	return nil
}

func (n *notify) NotifyAll(msg doki.IMessage) error {
	for id, v := range n.cimap {
		err := v.SendMsg(msg.GetMsgId(), msg.GetData())
		if err != nil {
			fmt.Println("Notify to", id, "err:", err)
			return err
		}
	}
	return nil
}

func (n *notify) NotifyBuffToConnByID(Id uint64, msg doki.IMessage) error {
	err := n.cimap[Id].SendBuffMsg(msg.GetMsgId(), msg.GetData())
	if err != nil {
		fmt.Println("Notify to", Id, "err:", err)
		return err
	}
	return nil
}

func (n *notify) NotifyBuffAll(msg doki.IMessage) error {
	for id, v := range n.cimap {
		err := v.SendBuffMsg(msg.GetMsgId(), msg.GetData())
		if err != nil {
			fmt.Println("Notify to", id, "err:", err)
			return err
		}
	}
	return nil
}
