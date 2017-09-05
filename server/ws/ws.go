package ws

import (
	"net"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/scgywx/foog"
)

type WebSocketServer struct{
	handle func(foog.IConn)
	msgType int
	upgrader websocket.Upgrader
}

type WebSocketConn struct{
	conn *websocket.Conn
	msgType int
	remoteAddr string
}

func NewServer()*WebSocketServer{
	s := &WebSocketServer{}
	s.upgrader.CheckOrigin = func(r *http.Request) bool{
		return true
	}
	return s
}

func (this *WebSocketServer)SetCheckOriginFunc(fn func(r *http.Request) bool){
	this.upgrader.CheckOrigin = fn
}

func (this *WebSocketServer)SetMessageType(msgType int){
	this.msgType = msgType
}

func (this *WebSocketServer)Run(ls net.Listener, fn func(foog.IConn)){
	this.handle = fn
	http.HandleFunc("/", this.handleConnection)
	http.Serve(ls, nil)
}

func (this *WebSocketServer)handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("websocket upgrade:", err)
		return
	}
	
	this.handle(&WebSocketConn{
		conn: c,
		msgType: this.msgType,
		remoteAddr: r.RemoteAddr,
	})
}

func (this *WebSocketConn)ReadMessage()([]byte, error){
	mt, msg, err := this.conn.ReadMessage()
	if err != nil{
		return nil, err
	}
	
	if this.msgType == 0 {
		this.msgType = mt
	}

	return msg, nil
}

func (this *WebSocketConn)WriteMessage(msg []byte) error{
	return this.conn.WriteMessage(this.msgType, msg)
}

func (this *WebSocketConn)Close(){
	this.conn.Close()
}

func (this *WebSocketConn)GetRemoteAddr()string{
	return this.remoteAddr
}