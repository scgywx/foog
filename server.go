package foog

import (
	"net"
	"fmt"
	"log"
	"time"
)

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

type servertask struct{
	sess *Session
	data []byte
}

var (
	wp *TaskPool
)

func RunServer(addr string, server IServer){
	ls, err := net.Listen("tcp", addr)
	if err != nil{
		fmt.Println("listen server failed", err)
		return 
	}

	log.Println("server started", addr)

	server.Run(ls, handleConnection)
}

func handleConnection(conn IConn){
	var (
		sess *Session
		err error
		msg []byte
	)

	sess = NewSession(conn)
	defer sess.Close()
	
	router.HandleAccept(sess)
	
	for{
		msg, err = sess.Conn.Recv()
		if err != nil{
			log.Println("read message error", err)
			break
		}
		
		err = handleRead(sess, msg)
		if err != nil{
			log.Println("handle message error", err)
		}
	}
	
	router.HandleClose(sess)
}

func handleRead(sess *Session, msg []byte) error{
	cmd, data, err := router.HandleRead(sess, msg)
	if err != nil{
		return err
	}

	err = Dispatch(cmd, &Context{
		Sess: sess,
		Data: data,
	})
	if err != nil{
		return err
	}

	sess.LastTime = time.Now().Unix()

	return nil
}