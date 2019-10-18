package foog

import (
	"os"
	"log"
	"github.com/scgywx/foog/utils"
)

type handlerFunc func(*Context)

var (
	nodeId int
	router IRouter
	logFile string
	logLevel int
	workerNum int
	handlers map[string]handlerFunc
	componentList []IComponent
)

func SetNodeId(id int){
	nodeId = id
}

func SetRouter(r IRouter){
	router = r
}

func SetWorkerNum(n int){
	workerNum = n
}

func SetLogLevel(level int){
	logLevel = level
}

func SetLogFile(filename string){
	logFile = filename
}

func Init(){
	var err error

	if err = initLog(); err != nil{
		log.Println("init log error", err)
		os.Exit(0)
	}

	initWorker()
	
	if err = initComponent(); err != nil{
		log.Println("init component error", err)
		os.Exit(0)
	}

	utils.SetUUIDNode(nodeId)
}

func init(){
	componentList = make([]IComponent, 0, 10)
	handlers = make(map[string]handlerFunc)
}