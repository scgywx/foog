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
		data []byte
	)

	sess = NewSession(conn)
	defer sess.Close()
	
	wp.Post(nil, func(interface{}) error{
		router.HandleAccept(sess)
		return nil
	})
	
	for{
		data, err = sess.Conn.Recv()
		if err != nil{
			log.Println("read message error", err)
			break
		}
		
		t := &servertask{
			sess: sess,
			data: data,
		}
		err = wp.Post(t, handleRead)
		if err != nil{
			log.Println("handle message error", err)
		}
	}
	
	wp.Post(nil, func(interface{}) error{
		router.HandleClose(sess)
		return nil
	})
}

func handleRead(arg interface{}) error{
	task := arg.(*servertask)
	cmd, data, err := router.HandleRead(task.sess, task.data)
	if err != nil{
		return err
	}

	err = Dispatch(cmd, &Context{
		Sess: task.sess,
		Data: data,
	})
	if err != nil{
		return err
	}

	task.sess.LastTime = time.Now().Unix()

	return nil
}

func initWorker(){
	wp = NewTaskPool(workerNum, 0)
	wp.Start()
}