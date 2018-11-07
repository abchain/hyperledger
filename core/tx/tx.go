package abchainTx

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/core/utils"
)

type tx struct {
	header proto.Message
	msgObj proto.Message
}

func msgToByte(m proto.Message) []byte {
	h, err := proto.Marshal(m)

	if err != nil {
		return nil
	}

	return h
}

func genHash(header []byte, msg []byte, method string) []byte {
	r, err := utils.DoubleSHA256(bytes.Join([][]byte{header, msg, []byte(method)}, nil))
	if err != nil {
		return nil
	}

	return r
}

func EncodeProto(m proto.Message) []byte {
	return msgToByte(m)
}

func DecodeProto(arg []byte, m proto.Message) error {
	return proto.Unmarshal(arg, m)
}

func (hasher *tx) GenHash(method string) []byte {

	hh := msgToByte(hasher.header)
	hm := msgToByte(hasher.msgObj)
	if hh == nil || hm == nil {
		return nil
	}

	return genHash(hh, hm, method)
}
