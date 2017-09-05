package foog

import(
	"fmt"
	"log"
	"reflect"
	"strings"
)

type IObject interface {
}

type handlerEntity struct {
	object IObject
	method reflect.Value
	argType reflect.Type
	argIsRaw bool
}

type handlerManager struct{
	handlers map[string]*handlerEntity
}

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
)

func isHandlerMethod(method reflect.Method) bool{
	mt := method.Type
	if mt.NumIn() != 3{
		return false
	}
	
	return true
}

func (this *handlerManager)register(obj IObject){
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	name := reflect.Indirect(v).Type().Name()
	
	if this.handlers == nil{
		this.handlers = make(map[string]*handlerEntity)
	}
	
	for m := 0; m < t.NumMethod(); m++ {
		method := t.Method(m)
		mt := method.Type
		mn := method.Name
		if isHandlerMethod(method){
			raw := false
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			
			this.handlers[strings.ToLower(fmt.Sprintf("%s.%s", name, mn))] = &handlerEntity{
				object: obj,
				method: v.Method(m),
				argType: mt.In(2),
				argIsRaw: raw,
			}
		}else{
			log.Printf("%s.%s register failed, argc=%d\n", name, mn, mt.NumIn())
		}
	}
}

func (this *handlerManager)dispatch(name string, sess *Session, data interface{}){
	h, ok := this.handlers[strings.ToLower(name)]
	if !ok {
		log.Println("not found handle by", name)
		return 
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("dispatch error", name, err)
		}
	}()
	
	var serialized bool
	var argv reflect.Value
	if !h.argIsRaw && sess.serializer != nil {
		if bytes,ok := data.([]byte); ok{
			argv = reflect.New(h.argType.Elem())
			err := sess.serializer.Decode(bytes, argv.Interface())
			if err != nil {
				log.Println("deserialize error", err.Error())
				return
			}

			serialized = true
		}
	}

	if !serialized {
		argv = reflect.ValueOf(data)
	}
	
	args := []reflect.Value{reflect.ValueOf(sess), argv}
	h.method.Call(args)
}