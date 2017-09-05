package foog

type ISerializer interface{
	Encode(interface{})([]byte, error)
	Decode([]byte, interface{})(error)
}