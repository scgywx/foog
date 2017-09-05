package foog

type groupEntity struct{
	name string
	members map[int64]*Session
}

type Group struct {
	groups map[string]*groupEntity
}

func NewGroup()*Group{
	group := &Group{
		groups: make(map[string]*groupEntity),
	}
	return group
}

func (this *Group)getGroup(name string)*groupEntity{
	grp, ok := this.groups[name]
	if !ok {
		grp = &groupEntity{
			name: name,
			members: make(map[int64]*Session),
		}
	}
	return grp
}

func (this *Group)Join(name string, sess *Session){
	grp := this.getGroup(name)
	grp.members[sess.Id] = sess
}

func (this *Group)Leave(name string, sess *Session){
	grp := this.getGroup(name)
	delete(grp.members, sess.Id)
}

func (this *Group)Clean(name string){
	delete(this.groups, name)
}

func (this *Group)Members(name string)([]*Session){
	members := []*Session{}
	grp, ok := this.groups[name]
	if ok{
		for _, s := range grp.members {
			members = append(members, s)
		}
	}
	return members
}

func (this *Group)Broadcast(name string, msg interface{}){
	grp, ok := this.groups[name]
	if !ok{
		return 
	}

	for _,v := range grp.members{
		v.WriteMessage(msg)
	}
}

func (this *Group)BroadcastWithoutSession(name string, msg interface{}, filters []*Session){
	grp, ok := this.groups[name]
	if !ok{
		return 
	}

	skip := make(map[int64]bool)
	for _,s := range filters{
		skip[s.Id] = true
	}

	for k,v := range grp.members{
		if _, ok := skip[k]; !ok{
			v.WriteMessage(msg)
		}
	}
}