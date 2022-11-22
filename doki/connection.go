package doki

import (
	"context"
	"errors"
	"fmt"
	"github.com/YanHeDoki/Doki/conf"
	"github.com/YanHeDoki/Doki/dokiIF"
	"github.com/YanHeDoki/Doki/pack"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"net"
	"sync"
	"time"
)

//当前链接模块
type Connection struct {
	//当前链接隶属于哪个Server
	TcpServer dokiIF.IServer
	//当前链接的 socket tcp 套接字
	Conn *net.TCPConn
	//链接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32

	//当前链接的状态
	IsClosed bool
	sync.RWMutex
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc

	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	MsgBuffChan chan []byte //定义channel成员

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

//NewConnection 创建连接的方法
func NewConnection(server dokiIF.IServer, conn *net.TCPConn, ConnID uint32) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      ConnID,
		IsClosed:    false,
		MsgBuffChan: make(chan []byte, conf.GlobalConfObject.MaxMsgChanLen), //不要忘记初始化
		property:    nil,
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {

	BaseLog.DefaultLog.DokiLog("debug", "Reader Server start ....")
	defer c.Stop()

	for {
		//检测是否关闭连接
		select {
		case <-c.ctx.Done():
			return
		default: //否则就普通业务操作
			message, err := c.TcpServer.GetPacket().UnPack(c.Conn)
			if err != nil {
				BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("UnPack err:%s", err))
				return
			}

			//得到当前数据的Request 数据
			req := Request{
				conn: c,
				msg:  message,
			}

			//已经设置开启了工作池
			if conf.GlobalConfObject.WorkerPoolSize > 0 {
				//发送消息到消息队列由工作池来处理
				c.TcpServer.GetMsgHandler().SendMsgToTaskQueue(&req)
			} else {
				//从路由中 找到注册绑定的Conn对应的router调用
				go c.TcpServer.GetMsgHandler().DoMsgHandler(&req)
			}
		}
	}
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWrite() {
	BaseLog.DefaultLog.DokiLog("debug", "Writer Goroutine is running")
	defer BaseLog.DefaultLog.DokiLog("debug", fmt.Sprint(c.RemoteAddr().String(), "conn Writer exit!"))

	for {
		select {
		//针对有缓冲channel需要些的数据处理
		case data, ok := <-c.MsgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Send Buff Data error:%s Conn Writer exit", err))
					return
				}
			} else {
				BaseLog.DefaultLog.DokiLog("warning", "MsgBuffChan is Closed")
				break
			}
		//用于退出
		case <-c.ctx.Done():
			//conn已经关闭
			return
		}
	}
}

//启动连接
func (c *Connection) Start() {
	BaseLog.DefaultLog.DokiLog("debug", fmt.Sprintf("conn starting...ConnID=%d", c.ConnID))
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//调用开发者设置的启动前的钩子函数
	c.TcpServer.CallOnConnStart(c)

	//启动当前链接的读数据业务
	go c.StartReader()
	// 启动当前链接的读数据业务
	go c.StartWrite()

	//阻塞
	select {
	case <-c.ctx.Done():
		//得到退出消息，不再阻塞
		c.finalizer()
		return
	}
}

func (c *Connection) Stop() {
	c.cancel()
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetContext() context.Context {
	return c.ctx
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}

	msg, err := c.TcpServer.GetPacket().Pack(pack.NewMsgPackage(msgId, data))
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Pack error msg ID =%d ", msgId))
		return errors.New("Pack error msg ")
	}

	//写回客户端
	_, err = c.Conn.Write(msg)
	return err
}

//带缓冲发送消息
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	//超市时间
	timeOut := time.NewTimer(10 * time.Millisecond)
	defer timeOut.Stop()

	if c.IsClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	msg, err := c.TcpServer.GetPacket().Pack(pack.NewMsgPackage(msgId, data))
	if err != nil {
		BaseLog.DefaultLog.DokiLog("error", fmt.Sprintf("Pack error msg ID =%d ", msgId))
		return errors.New("Pack error msg ")
	}

	// 发送超时
	select {
	case <-timeOut.C:
		return errors.New("send buff msg timeout")
	case c.MsgBuffChan <- msg:
		return nil
	}

}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	res, ok := c.property[key]
	if !ok {
		return nil, errors.New("not found key for Connection Property ")
	} else {
		return res, nil
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *Connection) finalizer() {

	//在销毁连接之前执行开发者的函数
	c.TcpServer.CallOnConnStop(c)
	//锁
	c.Lock()
	defer c.Unlock()
	//如果是已经关闭的就不用再处理
	if c.IsClosed {
		return
	}
	//关闭连接
	c.Conn.Close()
	//将当前链接从connmgr中销毁
	c.TcpServer.GetConnMgr().Remove(c)
	//关闭通道
	close(c.MsgBuffChan)
	//标志设置
	c.IsClosed = true

}
