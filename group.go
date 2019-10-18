package foog

import(
	"sync"
	"errors"
)

type Group struct {
	mutex sync.RWMutex
	once sync.Once
	members map[int64]*Session
}

var (
	errMemberNotFound = errors.New("member not found")
)

func NewGroup()*Group{
	group := &Group{
		members: make(map[int64]*Session),
	}
	return group
}

func (this *Group)Join(sess *Session){
	this.mutex.Lock()
	defer this.mutex.Unlock()
	
	this.members[sess.Id] = sess
}

func (this *Group)Leave(sess *Session){
	this.mutex.Lock()
	defer this.mutex.Unlock()
	
	delete(this.members, sess.Id)
}

func (this *Group)Clean(){
	this.mutex.Lock()
	defer this.mutex.Unlock()
	
	this.members = make(map[int64]*Session)
}

func (this *Group)Member(sid int64)(*Session, error){
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	
	mem, ok := this.members[sid]
	if ok {
		return mem, nil
	}
	
	return nil, errMemberNotFound
}

func (this *Group)Members()([]*Session){
	members := []*Session{}
	for _, v := range this.members {
		members = append(members, v)
	}
	return members
}

func (this *Group)Broadcast(msg interface{}){
	for _,v := range this.members{
		v.Send(msg)
	}
}

func (this *Group)BroadcastWithoutSession(msg interface{}, filters map[int64]bool){
	for k,v := range this.members{
		if _, ok := filters[k]; !ok{
			v.Send(msg)
		}
	}
}