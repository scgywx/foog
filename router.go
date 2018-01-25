package foog

type IRouter interface{
	HandleConnection(*Session)
	HandleClose(*Session)
	HandleMessage(*Session, interface{})(string, interface{}, error)
}
