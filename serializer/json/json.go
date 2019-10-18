package json

import (
	"encoding/json"
)

type JsonSerializer struct{}

func (this *JsonSerializer)Encode(v interface{})([]byte, error){
	return json.Marshal(v)
}

func (this *JsonSerializer)Decode(data []byte, v interface{}) error{
	return json.Unmarshal(data, v)
}