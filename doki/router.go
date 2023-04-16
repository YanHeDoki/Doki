package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	"strconv"
	"sync"
)

type Router struct {
	Apis     map[uint32][]dokiIF.RouterHandler
	Handlers []dokiIF.RouterHandler
	sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		Apis:     make(map[uint32][]dokiIF.RouterHandler, 10),
		Handlers: make([]dokiIF.RouterHandler, 0, 6),
	}
}

func (r *Router) Use(handles ...dokiIF.RouterHandler) {
	r.Handlers = append(r.Handlers, handles...)
}

func (r *Router) AddHandler(msgId uint32, Handlers ...dokiIF.RouterHandler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := r.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	finalSize := len(r.Handlers) + len(Handlers)
	mergedHandlers := make([]dokiIF.RouterHandler, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], Handlers)
	r.Apis[msgId] = append(r.Apis[msgId], mergedHandlers...)
}

func (r *Router) GetHandlers(MsgId uint32) ([]dokiIF.RouterHandler, bool) {
	r.RLock()
	defer r.RUnlock()
	handlers, ok := r.Apis[MsgId]
	return handlers, ok
}

func (r *Router) Group(start, end uint32, Handlers ...dokiIF.RouterHandler) dokiIF.IGroupRouter {
	return NewGroup(start, end, r, Handlers...)
}

type GroupRouter struct {
	start    uint32
	end      uint32
	Handlers []dokiIF.RouterHandler
	router   *Router
}

func NewGroup(start, end uint32, router *Router, Handlers ...dokiIF.RouterHandler) *GroupRouter {
	g := &GroupRouter{
		start:    start,
		end:      end,
		Handlers: make([]dokiIF.RouterHandler, 0, len(Handlers)),
		router:   router,
	}
	g.Handlers = append(g.Handlers, Handlers...)
	return g
}

func (g *GroupRouter) Use(Handlers ...dokiIF.RouterHandler) {
	g.Handlers = append(g.Handlers, Handlers...)
}

func (g *GroupRouter) AddHandler(MsgId uint32, Handlers ...dokiIF.RouterHandler) {
	if MsgId < g.start || MsgId > g.end {
		panic("add router to group err in msgId:" + strconv.Itoa(int(MsgId)))
	}

	finalSize := len(g.Handlers) + len(Handlers)
	mergedHandlers := make([]dokiIF.RouterHandler, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], Handlers)
	//回调实际路由的添加组件
	g.router.AddHandler(MsgId, mergedHandlers...)
}

// 以下均为废弃方法
//
//	func (r *Router) Next(request dokiIF.IRequest) {
//		r.Lock()
//		defer r.Unlock()
//		r.index++
//		for r.index < int8(len(r.Handlers)) {
//			r.Handlers[r.index](request)
//			r.index++
//		}
//		r.index = -1
//	}
//
//	func (r *Router) Abort() {
//		r.index = int8(len(r.Handlers))
//	}
//
//	func (r *Router) IsAbort() bool {
//		return r.index >= int8(len(r.Handlers))
//	}
//
//	func (r *Router) Reset() {
//		r.Lock()
//		defer r.Unlock()
//		r.index = -1
//		r.Handlers = make([]dokiIF.RouterHandler, 0, 1)
//	}
//
//	func (r *Router) Reindex() {
//		r.index = -1
//	}
