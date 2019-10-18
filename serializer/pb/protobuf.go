package pb

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

var errWrongValueType = errors.New("protobuf: convert on wrong type value")

type ProtobufSerializer struct{
}

func (s *ProtobufSerializer) Encode(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, errWrongValueType
	}
	return proto.Marshal(pb)
}

func (s *ProtobufSerializer) Decode(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errWrongValueType
	}
	return proto.Unmarshal(data, pb)
}