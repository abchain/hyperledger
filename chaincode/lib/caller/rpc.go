package rpc

import (
	_ "encoding/base64"
	"errors"
	"github.com/golang/protobuf/proto"
)

func EncodeRPCResult(msg proto.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("No result")
	}
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
