package foog

import(
	"log"
	"strings"
	"errors"
	"runtime/debug"
)

func Bind(cmd string, fn handlerFunc){
	handlers[strings.ToLower(cmd)] = fn
}

func UnBind(cmd string){
	delete(handlers, strings.ToLower(cmd))
}

func Dispatch(cmd string, ctx *Context)(err error){
	fn, ok := handlers[strings.ToLower(cmd)]
	if !ok {
		return errors.New("handle not found " + cmd)
	}
	
	defer func(){
		if r := recover(); r != nil{
			err = errors.New("handle dispatch error " + cmd)
			log.Println(r)
			log.Println(string(debug.Stack()))
		}
	}()

	fn(ctx)

	return 
}