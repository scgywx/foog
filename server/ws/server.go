package ws

import (
	"net"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/scgywx/foog"
)

type WSServer struct{
	handle func(foog.IConn)
	msgType int
	upgrader websocket.Upgrader
}

func NewServer()*WSServer{
	s := &WSServer{}
	s.upgrader.CheckOrigin = func(r *http.Request) bool{
		return true
	}
	return s
}

func (this *WSServer)SetCheckOriginFunc(fn func(r *http.Request) bool){
	this.upgrader.CheckOrigin = fn
}

func (this *WSServer)SetMessageType(msgType int){
	this.msgType = msgType
}

func (this *WSServer)Run(ls net.Listener, fn func(foog.IConn)){
	this.handle = fn
	http.HandleFunc("/", this.handleConnection)
	http.Serve(ls, nil)
}

func (this *WSServer)handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("websocket upgrade:", err)
		return
	}
	
	this.handle(&WSConn{
		server: this,
		conn: c,
		msgType: this.msgType,
		remoteAddr: r.RemoteAddr,
	})
}