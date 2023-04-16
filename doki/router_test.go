package doki

import (
	"fmt"
	"github.com/YanHeDoki/Doki/dokiIF"
	"testing"
)

func A1(request dokiIF.IRequest) {
	fmt.Println("test A1")
}
func A2(request dokiIF.IRequest) {
	fmt.Println("test A2")
}
func A3(request dokiIF.IRequest) {
	fmt.Println("test A3")
}
func A4(request dokiIF.IRequest) {
	fmt.Println("test A4")
}
func A5(request dokiIF.IRequest) {
	fmt.Println("test A5")
}

func TestRouterAdd(t *testing.T) {
	router := NewRouter()
	router.Use(A3)
	router.AddHandler(1, A1, A2)

	testgroup := router.Group(2, 5, A5)
	{
		testgroup.AddHandler(2, A4)

		//正确panic
		//testgroup.AddHandler(6, A4)
	}

	for _, v := range router.Apis[2] {
		v(&Request{
			conn:     nil,
			udpConn:  nil,
			msg:      nil,
			handlers: nil,
			index:    -1,
		})
	}
}
