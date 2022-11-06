package conf

import (
	"encoding/json"
	"github.com/YanHeDoki/Doki/iface"
	"io/ioutil"
	"os"
)

//存储全局参数 供其他模块使用
//一些参数可以通过配置文件由用户自定义

type GlobalObj struct {

	//Server
	TcpServer iface.IServer //当前全局的Server对象
	Host      string        //当前服务器主机监听的IP
	TcpPort   int           //当前服务器监听的端口
	Name      string        //当前服务器的名称

	//服务器可选配置
	Version          string //版本
	MaxConn          int    //最大连接数量
	MaxPacketSize    uint32 //当前框架数据包的最大尺寸
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
	PacketName       string //解包名称
	/*
		config file path
	*/
	ConfFilePath string
}

//定义一个全局的对外GlobalObj对象
var GlobalConfObject *GlobalObj

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (g *GlobalObj) Reload() {

	if confFileExists, _ := pathExists("./conf/Dokiconf.json"); !confFileExists {
		return
	}

	data, err := ioutil.ReadFile("./conf/Dokiconf.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
}

//提供一个Init方法初始化当前的全局对象
func ConfigInit() {
	//如果配置文件没有加载就是默认值
	GlobalConfObject = &GlobalObj{
		Name:             "ServerApp",
		Version:          "V1.0",
		Host:             "0.0.0.0",
		TcpPort:          9991,
		WorkerPoolSize:   10,
		MaxConn:          10000,
		MaxPacketSize:    4096,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    100,
		PacketName:       iface.StdDataPack,
	}

	//应该尝试从配置文件中的用户自定义的文件中读取
	GlobalConfObject.Reload()
}
