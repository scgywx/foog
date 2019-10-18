package tcp

import (
	"net"
	"io"
	"bufio"
	"encoding/binary"
	"github.com/scgywx/foog"
)

type TcpConn struct{
	server *TcpServer
	conn net.Conn
	br *bufio.Reader
	bw *bufio.Writer
	remoteAddr string
}

func (this *TcpConn)Recv()([]byte, error){
	headSize := this.server.headSize
	head := make([]byte, headSize)
	if _, err := io.ReadFull(this.br, head); err != nil{
		return nil, err
	}

	bodySize := 0
	if headSize == 2{
		bodySize = int(binary.BigEndian.Uint16(head))
	}else{
		bodySize = int(binary.BigEndian.Uint32(head) & 0x7fffffff)
	}

	body := make([]byte, bodySize)
	if _, err := io.ReadFull(this.br, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (this *TcpConn)Send(msg []byte)error{
	headSize := this.server.headSize
	bodySize := len(msg)
	hdr := make([]byte, headSize)

	if headSize == 2{
		binary.BigEndian.PutUint16(hdr, uint16(bodySize))
	}else{
		binary.BigEndian.PutUint32(hdr, uint32(bodySize))
	}

	if _, err := this.bw.Write(hdr); err != nil{
		return err
	}

	if _, err := this.bw.Write(msg); err != nil{
		return err
	}

	if err := this.bw.Flush(); err != nil{
		return err
	}
	
	return nil
}

func (this *TcpConn)Close(){
	this.conn.Close()
}

func (this *TcpConn)GetRemoteAddr()string{
	return this.remoteAddr
}

func (this *TcpConn)GetServer()foog.IServer{
	return this.server
}