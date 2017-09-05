package tcp

import (
	"net"
	"log"
	"bytes"
	"bufio"
	"encoding/binary"
	"github.com/scgywx/foog"
)

type TcpServer struct{
	headSize int
}

type TcpConn struct{
	server *TcpServer
	conn net.Conn
	br *bufio.Reader
	bw *bytes.Buffer
	remoteAddr string
}

func NewServer()*TcpServer{
	return &TcpServer{
		headSize: 4,
	}
}

func (this *TcpServer)SetHeadSize(n int){
	this.headSize = n
}

func (this *TcpServer)Run(ls net.Listener, fn func(foog.IConn)){
	if this.headSize != 2 && this.headSize != 4{
		this.headSize = 4
	}

	for {
		c, err := ls.Accept()
		if err != nil {
			log.Println("Accept failed", err)
			break
		}
		
		go fn(&TcpConn{
			server: this,
			conn: c,
			br: bufio.NewReaderSize(c, 1024),
			bw: bytes.NewBuffer(make([]byte, 0, 1024)),
			remoteAddr: c.RemoteAddr().String(),
		})
	}
}

func (this *TcpConn)ReadMessage()([]byte, error){
	headSize := this.server.headSize
	head, err := this.br.Peek(headSize)
	if err != nil{
		return nil, err
	}

	this.br.Discard(headSize)
	bodySize := 0
	if headSize == 2{
		bodySize = int(binary.BigEndian.Uint16(head))
	}else{
		bodySize = int(binary.BigEndian.Uint32(head) & 0x7fffffff)
	}

	off := 0
	bytes := make([]byte, bodySize)
	for off < bodySize{
		n, err := this.br.Read(bytes[off:])
		if err != nil{
			return  nil, err
		}

		off+= n
	}

	return bytes, nil
}

func (this *TcpConn)WriteMessage(msg []byte)error{
	headSize := this.server.headSize
	bodySize := len(msg)
	hdr := make([]byte, headSize)

	if headSize == 2{
		binary.BigEndian.PutUint16(hdr, uint16(bodySize))
	}else{
		binary.BigEndian.PutUint32(hdr, uint32(bodySize))
	}

	this.bw.Reset()
	this.bw.Write(hdr)
	this.bw.Write(msg)
	_, err := this.conn.Write(this.bw.Bytes())
	if err != nil{
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