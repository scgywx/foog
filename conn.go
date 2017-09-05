package foog

type IConn interface{
	ReadMessage()([]byte, error)
	WriteMessage([]byte) error
	Close()
	GetRemoteAddr()string
}