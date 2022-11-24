# Doki (心跳)



## 敬启者

##### Doki 是基于Zinx的二次开发的框架。

毕业设计选择是Tcp服务器，有过Golang服务器的实习经验，但没有底层实现的实践过，在Bilibii发现了zinx，在zinx 基础上按照自己偏好缝缝补补写了Doki作为服务端的框架来实现毕业设计，也希望能找到Dokidoki的工作。



## Doki改动了



##### 注入形式配置 （已PR并入zinx）

提供一个配置类，使用者可以根据具体的配置文件如Yaml Toml 或者是其他的配置一次性写入配置之后再使用  

```go
func NewServer(config *conf.Config) dokiIF.IServer 
```

函数启动服务器，其中传入的参数就是配置类，应当把配置解析到其中。

如果使用默认配置则是使用 

```go
func DefaultServer() dokiIF.IServer 
```



##### 路由操作

路由处理改动为函数集的形式，具体操作抽象成一个函数，router路由实现存储一个函数切片来存入每一个路由id对应的一个或一系列的操作

使用起来只需要实现一个或者多个如下函数即可

```go
func(request IRequest)
```

并且在request 里提供以下方法操作函数执行

```go
//执行下一个函数
Next(request IRequest)
//终结路由函数的执行
Abort()
//是否终结了函数
IsAbort() bool
```



##### 连接用户id映射

添加了一个Notify层用来做Id和用户链接的一个映射 ，通过框架统一来通知一些信息或者一次性通知所有人信息，也可以做到检测用户是否在线



##### 简单Udp 服务器

加入了简单的udp服务，主要配合notify来做到转发和同步信息使用





## 差异处使用



##### 快速启动：

配置好后 默认使用doki.DefaultServer可以启动默认服务器，如果要注入配置可以选择 doki.NewServer(conf) 方法启动服务器



##### 路由添加：



再doki中不在需要自己实现接口去做路由只需要满足类型的实现具体操作方法  如下

```go
func Handle1( req dokiIF.IRequest) {
   fmt.Println("1")
   if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
      fmt.Println(err)
   }
}

func Handle2( req dokiIF.IRequest) {
   req.Next()
   fmt.Println("2")
   if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
      fmt.Println(err)
   }
}

func Handle3( req dokiIF.IRequest) {
   fmt.Println("3")
   req.Abort()
   if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
      fmt.Println(err)
   }
}

func Handle4( req dokiIF.IRequest) {
   fmt.Println("4")
   if err := req.GetConnection().SendMsg(0, []byte("test4")); err != nil {
      fmt.Println(err)
   }

}
```

之后可以一次性添加到一个路由当中之后 再启服务，服务器就会自动构造路由并且添加操作函数

```go
s := doki.DefaultServer()
s.AddRouter(1, Handle1, Handle2, Handle3, Handle4)
s.Server()
```

并且IRequest接口提供的方法也可以直接使用生效  上列操作的结果是 打印出 1 3 2





##### Notify 通知操作

再具体项目或者使用中添加Notify操作的层 之后再实例化提供的接口就可以操作或者封装函数 如下，notify的方法均已上锁，可以直接使用

```go
var NT = Notify.NewNotify()

func AddUserConn(id uint64, conn dokiIF.IConnection) {
   NT.SetNotifyID(id, conn)
   conn.SetProperty("id", id)
}
func DelUserConn(id uint64) {
   NT.DelNotifyByID(id)
}

func SendUser(id uint64, MsgId uint32, data []byte) {
   err := NT.NotifyToConnByID(id, MsgId, data)
   if err != nil {
      fmt.Println(err)
      return
   }
}

func Broadcast(ids []uint64, MsgId uint32, data []byte) {
   for i := range ids {
      err := NT.NotifyToConnByID(ids[i], MsgId, data)
      if err != nil {
         fmt.Println(err)
         return
      }
   }
}
```









