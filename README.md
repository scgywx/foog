foog is a lightweight game server framework.

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