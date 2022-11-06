package doki

import (
	"github.com/YanHeDoki/Doki/iface"
)

type Router struct {
	index    int8 //函数索引
	handlers []iface.RouterHandler
}

func (r *Router) Next(request iface.IRequest) {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r, request)
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
	r.handlers = make([]iface.RouterHandler, 0, 1)
}

func (r *Router) Reindx() {
	r.index = -1
}
func (r *Router) AddHandler(handler ...iface.RouterHandler) {
	r.handlers = append(r.handlers, handler...)
}
