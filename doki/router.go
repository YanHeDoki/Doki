package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
	"sync"
)

type Router struct {
	index    int8 //函数索引
	handlers []dokiIF.RouterHandler
	sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		index:    -1,
		handlers: make([]dokiIF.RouterHandler, 0, 1),
	}
}

func (r *Router) Next(request dokiIF.IRequest) {
	r.Lock()
	defer r.Unlock()
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](request)
		r.index++
	}
	r.index = -1
}

func (r *Router) Abort() {
	r.index = int8(len(r.handlers))
}

func (r *Router) IsAbort() bool {
	return r.index >= int8(len(r.handlers))
}

func (r *Router) Reset() {
	r.Lock()
	defer r.Unlock()
	r.index = -1
	r.handlers = make([]dokiIF.RouterHandler, 0, 1)
}

func (r *Router) Reindex() {
	r.index = -1
}
func (r *Router) AddHandler(handler ...dokiIF.RouterHandler) {
	r.handlers = append(r.handlers, handler...)
}
