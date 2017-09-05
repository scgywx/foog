package foog

import (
	"fmt"
	"net"
	"log"
	"os"
	"time"
	"sync"
)

type IServer interface{
	Run(net.Listener, func(IConn))
}

type IConn interface{
	ReadMessage()([]byte, error)
	WriteMessage([]byte) error
	Close()
	GetRemoteAddr()string
}

type IRouter interface{
	HandleConnection(*Session)
	HandleClose(*Session)
	HandleMessage(*Session, []byte)(string, interface{}, error)
}

type Application struct{
	id int
	listenAddr string
	server IServer
	router IRouter
	serializer ISerializer
	logFile string
	logLevel int
	handler handlerManager
	sessions map[int64]*Session
	mutexSess sync.Mutex
}

func (this *Application)Register(c IObject){
	this.handler.register(c)
}

func (this *Application)Listen(addr string){
	this.listenAddr = addr
}

func (this *Application)SetServer(s IServer){
	this.server = s
}

func (this *Application)SetRouter(r IRouter){
	this.router = r
}

func (this *Application)SetLogLevel(level int){
	this.logLevel = level
}

func (this *Application)SetLogFile(filename string){
	this.logFile = filename
}

func (this *Application)SetSerializer(s ISerializer){
	this.serializer = s
}

func (this *Application)SetId(id int){
	this.id = id
}

func (this *Application)GetSessionById(sid int64)*Session{
	return this.sessions[sid]
}

func (this *Application)Start(){
	//init server
	ls, err := net.Listen("tcp", this.listenAddr)
	if err != nil{
		fmt.Println("listen server failed", err)
		return 
	}

	//init log
	if len(this.logFile) > 0 {
		w, err := os.OpenFile(this.logFile, os.O_RDWR | os.O_APPEND | os.O_CREATE, os.ModePerm)
		if err != nil{
			fmt.Println("open log file error", err)
			return 
		}
		
		log.SetOutput(w)
	}
	
	this.sessions = make(map[int64]*Session)
	
	log.Println("server started", this.listenAddr)
	this.server.Run(ls, this.handleConnection)
}

func (this *Application)handleConnection(conn IConn){
	sess := NewSession(conn, this.id)
	sess.serializer = this.serializer

	this.mutexSess.Lock()
	this.sessions[sess.Id] = sess	
	this.mutexSess.Unlock()

	defer func(){
		this.mutexSess.Lock()
		delete(this.sessions, sess.Id)
		this.mutexSess.Unlock()
		conn.Close()
	}()
	defer this.router.HandleClose(sess)
	
	this.router.HandleConnection(sess)

	for{
		msg, err := conn.ReadMessage()
		if err != nil{
			log.Println("read message failed", err)
			break
		}
		
		name, data, err := this.router.HandleMessage(sess, msg)
		if err != nil{
			log.Println("handle message failed", err)
			break
		}
		
		sess.LastTime = time.Now().Unix()
		this.handler.dispatch(name, sess, data)
	}
}
