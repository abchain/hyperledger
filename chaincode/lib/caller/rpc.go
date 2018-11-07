package rpc

import (
	_ "encoding/base64"
	"github.com/golang/protobuf/proto"
)

func EncodeRPCResult(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func DecodeRPCResult(msg proto.Message, r []byte) error {
	return proto.Unmarshal(r, msg)
}

type Caller interface {
	Deploy(method string, arg [][]byte) (string, error)
	Invoke(method string, arg [][]byte) (string, error)
	Query(method string, arg [][]byte) ([]byte, error)
}
