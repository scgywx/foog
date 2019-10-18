package foog

type IRouter interface{
	HandleAccept(sess *Session)
	HandleRead(sess *Session, data []byte)(string, interface{}, error)
	HandleWrite(sess *Session, msg interface{})([]byte, error)
	HandleClose(sess *Session)
}