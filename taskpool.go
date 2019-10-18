package foog

import (
	"log"
)

type taskFn func(interface{})error

type taskEntity struct{
	fn taskFn
	data interface{}
	done chan error
}

type TaskPool struct{
	workerNum int
	queue chan *taskEntity
}

func NewTaskPool(workerNum int, queueNum int)*TaskPool{
	tp := &TaskPool{
		workerNum: workerNum,
		queue : make(chan *taskEntity, queueNum),
	}
	return tp
}

func (this *TaskPool)Start(){
	for i := 0; i < this.workerNum; i++{
		go this.runWorker(i+1)
	}
}

func (this *TaskPool)AsyncPost(data interface{}, fn taskFn){
	this.postRaw(data, nil, fn)
}

func (this *TaskPool)Post(data interface{}, fn taskFn)error{
	return this.postRaw(data, make(chan error), fn)
}

func (this *TaskPool)postRaw(data interface{}, done chan error, fn taskFn)error{
	var err error

	if workerNum > 0{
		this.queue <- &taskEntity{
			fn: fn,
			data: data,
			done: done,
		}

		if done != nil{
			err = <-done
		}
	}else{
		err = fn(data)
		log.Println("runat current co")
	}

	return err
}

func (this *TaskPool)runWorker(i int){
	for{
		task := <-this.queue
		err := task.fn(task.data)
		if task.done != nil{
			task.done <- err
		}
		log.Println("runat taskpool co", i)
	}
}