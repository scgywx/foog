package foog

type IRouter interface{
	HandleConnection(*Session)
	HandleClose(*Session)
	HandleMessage(*Session, []byte)(string, interface{}, error)
}
