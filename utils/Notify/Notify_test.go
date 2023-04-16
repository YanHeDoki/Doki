package Notify

import (
	"fmt"
	"github.com/YanHeDoki/Doki/doki"
	"github.com/YanHeDoki/Doki/dokiIF"
	"github.com/YanHeDoki/Doki/pack"
	"net"
	"strconv"
	"testing"
	"time"
)

var nt = NewNotify()

func Handle1(req dokiIF.IRequest) {
	id, _ := strconv.Atoi(string(req.GetData()))
	nt.SetNotifyID(uint64(id), req.GetConnection())
}

func Server() {
	s := doki.DefaultServer()
	s.AddRouter(1, Handle1)
	s.Server()
}

func Clinet() {
	//conf.ConfigInit()
	//1创建直接链接
	for i := 0; i < 7000; i++ {
		go func(i int) {
			conn, err := net.Dial("tcp", "127.0.0.1:9512")
			if err != nil {
				fmt.Println("net dial err:", err)
				return
			}
			defer conn.Close()
			//链接调用write方法写入数据
			id := strconv.Itoa(i)
			dp := pack.NewDataPack()
			msg, err := dp.Pack(pack.NewMsgPackage(1, []byte(id)))
			if err != nil {
				return
			}
			_, err = conn.Write(msg)

			if err != nil {
				return
			}
			select {}
			//fmt.Println("==> Recv Msg: ID=", NewMsg.GetMsgId(), ", len=", NewMsg.GetDataLen(), ", data=", string(NewMsg.GetData()))
		}(i)
		//time.Sleep(1 * time.Nanosecond)
	}
}

func init() {
	go Server()
	go Clinet()
	//go ClinetJoin()
}

func ClinetJoin() {
	t := time.NewTicker(500 * time.Millisecond)
	i := 10000
	for {
		select {
		case <-t.C:
			go func(i int) {
				conn, err := net.Dial("tcp", "127.0.0.1:9512")
				if err != nil {
					fmt.Println("net dial err:", err)
					return
				}
				defer conn.Close()
				//链接调用write方法写入数据
				id := strconv.Itoa(i)
				dp := pack.NewDataPack()
				msg, err := dp.Pack(pack.NewMsgPackage(1, []byte(id)))
				if err != nil {
					return
				}
				_, err = conn.Write(msg)

				if err != nil {
					return
				}
				select {}
				//fmt.Println("==> Recv Msg: ID=", NewMsg.GetMsgId(), ", len=", NewMsg.GetDataLen(), ", data=", string(NewMsg.GetData()))
			}(i)
			i++
		}
	}

}

func TestAA(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println(len(nt.cimap))
	})
	time.Sleep(6 * time.Second)
}

func BenchmarkNotify(b *testing.B) {
	time.Sleep(6 * time.Second)
	fmt.Println(len(nt.cimap))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nt.NotifyAll(1, []byte("雪下的是盐"))
	}
}
