package doki

//尝试修改中
type RouterHandler func(router IRouter, request IRequest)
type Router struct {
	index    int8 //函数索引
	handlers []RouterHandler
}

func (r *Router) Next(request IRequest) {
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
	r.handlers = make([]RouterHandler, 0, 1)
}

func (r *Router) Reindx() {
	r.index = -1
}
func (r *Router) AddHandler(handler ...RouterHandler) {
	r.handlers = append(r.handlers, handler...)
}
