package conf

import (
	"encoding/json"
	"github.com/YanHeDoki/Doki/constants"
	"github.com/YanHeDoki/Doki/dokiIF"
	"io/ioutil"
	"os"
	"runtime"
)

//存储全局参数 供其他模块使用
//一些参数可以通过配置文件由用户自定义

type GlobalObj struct {

	//Server
	TcpServer dokiIF.IServer //当前全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器监听的端口
	Name      string         //当前服务器的名称

	//服务器可选配置
	Version          string //版本
	MaxConn          int    //最大连接数量
	MaxPacketSize    uint32 //当前框架数据包的最大尺寸
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	DoMsgHandlerNum  int    //一个消息池子多少个线程去执行任务
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
	PacketName       string //解包名称
	/*
		config file path
	*/
	ConfFilePath string
	Loglevel     string
}

// 定义一个全局的对外GlobalObj对象
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

// 提供一个Init方法初始化当前的全局对象
func ConfigInit() {
	WorkerPoolSize := runtime.NumCPU()
	//如果配置文件没有加载就是默认值
	GlobalConfObject = &GlobalObj{
		Name:             "ServerApp",
		Version:          "V1.0",
		Host:             "127.0.0.1",
		TcpPort:          9512,
		WorkerPoolSize:   uint32(WorkerPoolSize),
		DoMsgHandlerNum:  3,
		MaxConn:          1000000,
		MaxPacketSize:    4096,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    100,
		PacketName:       constants.StdDataPack,
		Loglevel:         "info",
	}

	//应该尝试从配置文件中的用户自定义的文件中读取
	GlobalConfObject.Reload()
}

func UserConfInit(config *Config) {

	//先使用默认配置
	ConfigInit()

	//如果配置文件没有加载就是默认值
	// Server
	if config.Name != "" {
		GlobalConfObject.Name = config.Name
	}
	if config.Host != "" {
		GlobalConfObject.Host = config.Host
	}
	if config.TcpPort != 0 {
		GlobalConfObject.TcpPort = config.TcpPort
	}
	if config.Version != "" {
		GlobalConfObject.Version = config.Version
	}

	//选择配置
	if config.MaxPacketSize != 0 {
		GlobalConfObject.MaxPacketSize = config.MaxPacketSize
	}
	if config.MaxConn != 0 {
		GlobalConfObject.MaxConn = config.MaxConn
	}
	if config.WorkerPoolSize != 0 {
		GlobalConfObject.WorkerPoolSize = config.WorkerPoolSize
	}
	if config.MaxWorkerTaskLen != 0 {
		GlobalConfObject.MaxWorkerTaskLen = config.MaxWorkerTaskLen
	}
	if config.MaxMsgChanLen != 0 {
		GlobalConfObject.MaxMsgChanLen = config.MaxMsgChanLen
	}
	if config.DoMsgHandlerNum != 0 {
		GlobalConfObject.DoMsgHandlerNum = config.DoMsgHandlerNum
	}
	if config.LogLevel != "" {
		GlobalConfObject.Loglevel = config.LogLevel
	}
}

// 注意如果使用UserConf应该调用方法同步至 GlobalConfObject 因为其他参数是调用的此结构体参数
func UserConfToGlobal(config *Config) {

}
