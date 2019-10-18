package tcp

import (
	"net"
	"bufio"
	"log"
	"github.com/scgywx/foog"
)

type TcpServer struct{
	headSize int
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
			br: bufio.NewReader(c),
			bw: bufio.NewWriter(c),
			remoteAddr: c.RemoteAddr().String(),
		})
	}
}