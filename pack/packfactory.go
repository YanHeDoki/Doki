package pack

import (
	"github.com/YanHeDoki/Doki/iface"
	"sync"
)

var pack_once sync.Once

type pack_factory struct{}

var factoryInstance *pack_factory

/*
	生成不同封包解包的方式，单例
*/
func Factory() *pack_factory {
	pack_once.Do(func() {
		factoryInstance = new(pack_factory)
	})

	return factoryInstance
}

//NewPack 创建一个具体的拆包解包对象
func (f *pack_factory) NewPack(kind string) iface.IDataPack {
	var dataPack iface.IDataPack

	switch kind {
	// 标准默认封包拆包方式
	case iface.StdDataPack:
		dataPack = NewDataPack()

		//case 自定义封包拆包方式case

	default:
		dataPack = NewDataPack()
	}
	return dataPack
}
