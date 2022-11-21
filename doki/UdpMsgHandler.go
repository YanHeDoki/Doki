package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	BaseLog "github.com/YanHeDoki/Doki/utils/log"
	"strconv"
)

type UdpMsgHandle struct {
	Apis map[uint32]*Router //路由模块
}

func NewUdpMsgHandle() *UdpMsgHandle {
	return &UdpMsgHandle{
		Apis: make(map[uint32]*Router),
	}
}

func (u *UdpMsgHandle) DoMsgHandler(request dokiIF.IRequest) {
	router, ok := u.Apis[request.GetMsgId()]
	if !ok {
		BaseLog.DefaultLog.DokiLog("warning", "not find Router In Apis")
		return
	}
	request.BindRouter(router)
	router.Reindx()
	router.Next(request)
}

func (u *UdpMsgHandle) AddRouter(msgId uint32, handler ...dokiIF.RouterHandler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := u.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	//优化一次完成再存入map
	r := NewRouter()
	r.AddHandler(handler...)
	u.Apis[msgId] = r
}
