package ws

import (
	"github.com/gorilla/websocket"
	"github.com/scgywx/foog"
)

type WSConn struct{
	server *WSServer
	conn *websocket.Conn
	msgType int
	remoteAddr string
}

func (this *WSConn)Recv()([]byte, error){
	mt, msg, err := this.conn.ReadMessage()
	if err != nil{
		return nil, err
	}
	
	if this.msgType == 0 {
		this.msgType = mt
	}

	return msg, nil
}

func (this *WSConn)Send(msg []byte) error{
	return this.conn.WriteMessage(this.msgType, msg)
}

func (this *WSConn)Close(){
	this.conn.Close()
}

func (this *WSConn)GetRemoteAddr()string{
	return this.remoteAddr
}

func (this *WSConn)GetServer()foog.IServer{
	return this.server
}