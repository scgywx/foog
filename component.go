package foog

type IComponent interface {
	Init()error
	Shutdown()
}

func initComponent()error{
	var err error
	for _,c := range componentList{
		if err = c.Init(); err != nil{
			return err
		}
	}
	return nil
}

func AddComponent(c IComponent){
	componentList = append(componentList, c)
}

func ShutdownComponent(){
	for _,c := range componentList{
		c.Shutdown()
	}
}