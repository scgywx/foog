foog是一个易使用、易扩展、轻量级的服务端框架.  
易扩展的协议支持，内置tcp(2 or 4bytes消息头+消息体)和websocket，仅需一行代码即可轻松实现协议切换，也可以同时支持两种协议，如果内置协议不满足需求，仅需要实现指定interface即可完成协议自定义.  
更自由的路由分发，根据不同的业务需求打造属于自己的路由规则，让一切动向全掌握在自己手中.   

## Installation
```shell
go get github.com/scgywx/foog
```

## Server
```go
package main

import (
	"fmt"
	"encoding/json"
	"github.com/scgywx/foog"
	"github.com/scgywx/foog/server/ws"
)

type MyRequest struct{
	Cmd string `json:"cmd"`
	Data map[string]string `json:"data"`
}

type MyResponse struct{
	Cmd string `json:"cmd"`
	Data interface{} `json:"data"`
}

type MyRouter struct{
}

func (this *MyRouter)HandleAccept(sess *foog.Session){
	fmt.Printf("new client from %s, #%d\n", sess.Conn.GetRemoteAddr(), sess.Id)
}

func (this *MyRouter)HandleClose(sess *foog.Session){
	fmt.Printf("client close #%d\n", sess.Id)
}

func (this *MyRouter)HandleRead(sess *foog.Session, msg []byte)(string, interface{}, error){
	req := &MyRequest{}
	json.Unmarshal(msg, req)
	return req.Cmd, req, nil
}

func (this *MyRouter)HandleWrite(sess *foog.Session, msg interface{})([]byte, error){
	data, err := json.Marshal(msg)
	return data, err
}

type SayResult struct{
	Text string `json:"text"`
}

func handle_Hello_Say(ctx *foog.Context){
	req := ctx.Data.(*MyRequest)
	res := &MyResponse{
		Cmd: req.Cmd,
		Data: &SayResult{
			Text: fmt.Sprintf("hello %s", req.Data["name"]),
		},
	}
	ctx.Sess.Send(res)
}

func main() {
	foog.SetNodeId(1)    //设置节点ID，用于UUID生成
	foog.SetWorkerNum(2) //设置最大工作线程
	foog.SetRouter(&MyRouter{}) //自定义Router
	foog.Bind("Hello.Say", handle_Hello_Say) 
	foog.Init()  //初始化框架
	foog.RunServer(":8888", ws.NewServer())//启动服务
}
```

## Client
```js
ws = new window.WebSocket("ws://127.0.0.1:8888")
ws.onmessage = function(v){
	console.log(v.data)
}
ws.send('{"cmd":"Hello.Say","data":{"name":"foog"}}')
```

## Router
每个应用必须设置一个Router用来处理请求分发，Router需要实现IRouter的3个方法
```go
type IRouter interface{
	HandleAccept(*Session)
	HandleClose(*Session)
	HandleRead(*Session, []byte)(string, interface{}, error)
	HandleWrite(*Session, interface{})([]byte, error)
}
```
HandleAccept和HandleClose分别在新连接和关闭连接时会调用，参数仅有一个Session.  
HandleRead在收到消息会调用，用于消息在应用层解包
HnadleWrite在要发送消息时调用，用于消息在应用层打包

## Protocol
内置tcp和websocket协议，而tcp采用消息头(2或者4字节)+消息体，消息头均使用大端模式，修改长度可用如下方式：
```go
svr := tcp.NewServer()
svr.SetHeadSize(2)//设置为2字节
```

ws协议默认会允许所有客户端连接，如果需要修改检查来源函数，代码如下：
```go
svr := ws.NewServer()
svr.SetCheckOriginFunc(func(r *http.Request) bool{
	//TODO
})
```

如果内置协议不满足需求，只需实现IServer和IConn接口即可
```go
type IServer interface{
	Run(net.Listener, func(IConn))
}
type IConn interface{
	Recv()([]byte, error)
	Send([]byte) error
	Close()
	GetRemoteAddr()string
	GetServer()IServer
}
```

应用已经监听了端口，IServer.Run实现部分仅需要accept，然后将接受到的连接以回调的方式传递给指定函数.  
当应用收到有新的连接，会直接调用IConn.ReadMessage，该方法直到收到一条完整的消息才会返回，那么协议、缓存、分包都需要在该方法内完成.  
当应用需要发消息就会调用IConn.Send，该方法需要实现通信协议打包、发送工作.  
IConn.Close实现关闭连接.  
IConn.GetRemoteAddr实现获取客户端IP和端口.  
IConn.GetServer返回当前服务器指针
具体的实现方式可参见server/tcp/server.go 与 server/tcp/conn.go