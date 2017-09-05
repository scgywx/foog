package foog

import(
	"time"
)

type Session struct {
	serializer ISerializer
	Id int64
	Conn IConn
	LastTime int64
	Data interface{}
}

var counter int64 = 0

/**
 * 1位符号
 * 31位时间戳(最大可表示到2038年)
 * 10位毫秒
 * 10位服务器ID(最大可表示1024)
 * 12位自增id(最大值是4096)
 * 共64位，每秒可生成400w条不同ID
 */
func NewSession(conn IConn, appId int)*Session{
	counter++
	sess := &Session{
		Id: ((time.Now().UnixNano() / 1000000) << 22) | int64((appId & 0x3ff) << 12) | (counter & 0xfff),
		Conn: conn,
	}
	return sess
}

func (this *Session)WriteMessage(data interface{}) error{
	if msg, ok := data.([]byte); ok || this.serializer == nil{
		return this.Conn.WriteMessage(msg)
	}else{
		bytes, err := this.serializer.Encode(data)
		if err != nil{
			return err
		}

		return this.Conn.WriteMessage(bytes)
	}
}