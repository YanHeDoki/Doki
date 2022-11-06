package Notify

import (
	"fmt"
	"github.com/YanHeDoki/Doki/iface"
	"sync"
)

type notify struct {
	cimap ConnIDMap
	look  sync.RWMutex
}

func NewNotify() *notify {
	return &notify{
		cimap: make(map[uint64]iface.IConnection, 100),
	}
}

func (n *notify) SetNotifyID(Id uint64, conn iface.IConnection) {
	n.look.Lock()
	defer n.look.RLock()
	n.cimap[Id] = conn
}
func (n *notify) GetNotifyByID(Id uint64) iface.IConnection {
	n.look.RLock()
	defer n.look.RLock()
	return n.cimap[Id]
}

func (n *notify) NotifyToConnByID(Id uint64, msg iface.IMessage) error {
	err := n.cimap[Id].SendMsg(msg.GetMsgId(), msg.GetData())
	if err != nil {
		fmt.Println("Notify to", Id, "err:", err)
		return err
	}
	return nil
}

func (n *notify) NotifyAll(msg iface.IMessage) error {
	for id, v := range n.cimap {
		err := v.SendMsg(msg.GetMsgId(), msg.GetData())
		if err != nil {
			fmt.Println("Notify to", id, "err:", err)
			return err
		}
	}
	return nil
}
