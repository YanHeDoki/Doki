package doki

import (
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/conf"
	"github.com/YanHeDoki/Doki/iface"
	"sync"
)

//链接管理模块

type ConnManager struct {
	connections map[uint32]iface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection, conf.GlobalConfObject.MaxConn/2),
	}

}

func (c *ConnManager) Add(connection iface.IConnection) {
	//保护共享资源 加锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[connection.GetConnID()] = connection
	fmt.Println("ADD conn to manager success")
}

func (c *ConnManager) Remove(connection iface.IConnection) {
	//保护共享资源 加锁
	c.connLock.Lock()
	delete(c.connections, connection.GetConnID())
	c.connLock.Unlock()
	fmt.Println("Remove conn to manager success")
}

func (c *ConnManager) Get(connId uint32) (iface.IConnection, error) {
	//保护共享资源 加锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (c *ConnManager) Len() int {
	c.connLock.RLock()
	length := len(c.connections)
	c.connLock.RUnlock()
	return length
}

//停止所有链接并逐个清理
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connId, conn := range c.connections {
		//停止这个连接的资源
		conn.Stop()
		delete(c.connections, connId)
	}
	fmt.Println("clear ConnManagerMap success")
}

//ClearOneConn  利用ConnID获取一个链接 并且删除
func (c *ConnManager) ClearOneConn(connID uint32) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	connections := c.connections
	if conn, ok := connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}
