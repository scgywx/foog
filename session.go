package foog

import(
	"log"
	"github.com/scgywx/foog/utils"
)

type Session struct {
	Id int64
	Conn IConn
	LastTime int64
	Closed bool
	Data interface{}
}

func NewSession(conn IConn)*Session{
	sess := &Session{
		Id: utils.UUID(),
		Conn: conn,
	}
	return sess
}

func (this *Session)Send(res interface{}) error{
	data, err := router.HandleWrite(this, res)
	if err != nil{
		log.Println("router handle send failed", err)
		return err
	}

	return this.Conn.Send(data)
}

func (this *Session)Close(){
	this.Closed = true
	this.Conn.Close()
}

