foog是一个易使用、易扩展、轻量级的服务端框架，foog旨为解决基础建设，所以并没有包含其它业务相关的辅助库，这也使得他更为精简.  
易扩展的协议支持，内置tcp(2或者4bytes消息头+消息体)和websocket，仅需一行代码即可轻松实现协议切换，如果内置协议不满足需求，仅需要实现指定interface即可完成协议自定义.  
更自由的路由分发，根据不同的业务需求，设计属于自己路由规则，让一切动向全掌握在自己手中.  
可配置的序列化方式，json、pb你想用，随时可以自行配置，如不满足其需求，自行实现Encode和Decode后配置即可.  

## Installation
```shell
go get github.com/scgywx/foog
```

## Hello-Server(using websocket protocol)
```go
package main

import (
	"fmt"
	"encoding/json"
	"github.com/scgywx/foog"
	"github.com/scgywx/foog/server/ws"
	sjson "github.com/scgywx/foog/serializer/json"
)

type MyRequest struct{
	Cmd string `json:"cmd"`
	Data map[string]interface{} `json:"data"`
}

type MyRouter struct{
}

func (this *MyRouter)HandleConnection(sess *foog.Session){
	fmt.Printf("new client from %s, #%d\n", sess.Conn.GetRemoteAddr(), sess.Id)
}

func (this *MyRouter)HandleClose(sess *foog.Session){
	fmt.Printf("client close #%d\n", sess.Id)
}

func (this *MyRouter)HandleMessage(sess *foog.Session, msg []byte)(string, interface{}, error){
	req := &MyRequest{}
	json.Unmarshal(msg, req)
	return req.Cmd, req.Data, nil
}

type Hello struct{
}

type SayResponse struct{
	Name string `json:"name"`
	Text string `json:"text"`
}

func (this *Hello)Say(sess *foog.Session, req map[string]interface{}){
	rsp := &SayResponse{
		Name: req["name"].(string),
		Text: fmt.Sprintf("hello %s", req["name"]),
	}
	fmt.Println(rsp.Text)
	sess.WriteMessage(rsp)
}

func main() {
	app := &foog.Application{}
	app.SetRouter(&MyRouter{})
	app.SetServer(ws.NewServer())
	app.SetSerializer(sjson.New())
	app.Register(&Hello{})
	app.Listen("127.0.0.1:9005")
	app.Start()
}
```

## Hello-Client
```js
ws = new window.WebSocket("ws://127.0.0.1:9005")
ws.onmessage = function(v){
	console.log(v)
}
ws.send('{"cmd":"Hello.Say","data":{"name":"test"}}')
```

## Application
每个应用需要创建一个实例，然后设置相关的参数，最后启动即可，如：
```go
app := &foog.Application{}//创建App
app.SetRouter(&MyRouter{})//设置Router [必选]
app.SetServer(ws.NewServer())//设置服务端通信协议为websocket,想改tcp就直接将ws改成tcp即可 [必选]
app.SetSerializer(sjson.New())//设置序列化方式，[可选]
app.Register(&Hello{})//注册模块 [必选]
app.Listen("127.0.0.1:9005")//监听端口 [必选]
app.Start()//启动
```

如果想启动多个应用，就实例化多个application，还可以实现不同端口使用不同协议，逻辑不变，手机使用tcp,网页使用websocket.

## Protocol
内置tcp和websocket协议，而tcp采用消息头(2或者4字节)+消息体，消息头均使用大端模式，修改长度可用如下方式：
```go
svr := tcp.NewServer()
svr.SetHeadSize(2)//设置为2字节
```

ws协议如果需要修改请求源检查函数，代码如下：
```go
svr := ws.NewServer()
svr.SetCheckOriginFunc(func(r *http.Request) bool{
	//TODO
})
```

如果内置协议不满足需求，那么需要实现server.go和conn.go提供的IServer和IConn接口
```go
type IServer interface{
	Run(net.Listener, func(IConn))
}
type IConn interface{
	ReadMessage()([]byte, error)
	WriteMessage([]byte) error
	Close()
	GetRemoteAddr()string
}
```

应用已经监听了配置的端口，那么在IServer.Run实现部分仅需要accept，然后将接受到的连接以回调的方式传递给指定函数.  
当应用收到有新的连接，会直接调用IConn.ReadMessage，该方法直到收到一条完整的消息才会返回，那么协议、缓存、分包都需要在该方法内完成.  
当应用需要发消息就会调用IConn.WriteMessage，该方法需要实现协议打包、发送工作.  
IConn.Close实现关闭连接.  
IConn.GetRemoteAddr实现获取客户端IP和端口.  
具体的实现方式可参见server/tcp/server.go.  

## Router
每个应用必须设置一个Router用来处理请求分发，Router需要实现：HandleConnection、HandleClose、HandleMessage三个方法.  
HandleConnection和HandleClose在新连接和关闭连接时会调用，参数仅有一个Session.  
HandleMessage在收到消息会调用，方法有两个参数：session和data，返回三个参数：handleName, packet, error.  
handleName即在app.Register的时候注册的结构体对象，其规则是：小写(结构体名.方法名)，如Hello的Say方法，注册后的handleName为hello.say.  
packet即本次传给指定模块方法的参数，在满足三个条件的情况下，将会自动转换packet类型，1、设置了消息序列化；2、该参数是[]byte类型；3、目标方法第二个参数非[]byte类型.  
error则表示由路出错，不予分发.  