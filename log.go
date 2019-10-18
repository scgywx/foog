package foog

import (
	"log"
	"os"
)

func initLog() error{
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if len(logFile) > 0 {
		w, err := os.OpenFile(logFile, os.O_RDWR | os.O_APPEND | os.O_CREATE, os.ModePerm)
		if err != nil{
			return err
		}
		
		log.SetOutput(w)
	}

	return nil
}

func Log(level int, _fmt string, _args ...interface{}){
	if level < logLevel{
		return 
	}

	log.Printf(_fmt, _args...)
}