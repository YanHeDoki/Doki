package doki

import (
	"github.com/YanHeDoki/Doki/dokiIF"
)

type Router struct {
	index    int8 //函数索引
	handlers []dokiIF.RouterHandler
}

func NewRouter() *Router {
	return &Router{
		index:    -1,
		handlers: make([]dokiIF.RouterHandler, 0, 1),
	}
}

func (r *Router) Next(request dokiIF.IRequest) {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](request)
		r.index++
	}
}

func (r *Router) Abort() {
	r.index = int8(len(r.handlers))
}

func (r *Router) IsAbort() bool {
	return r.index >= int8(len(r.handlers))
}

func (r *Router) Reset() {
	r.index = -1
	r.handlers = make([]dokiIF.RouterHandler, 0, 1)
}

func (r *Router) Reindx() {
	r.index = -1
}
func (r *Router) AddHandler(handler ...dokiIF.RouterHandler) {
	r.handlers = append(r.handlers, handler...)
}
