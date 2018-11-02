package abchainTx

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/core/utils"
)

type tx struct {
	header proto.Message
	msgObj proto.Message
}

func toArgument(b []byte) string {
	return base64.RawStdEncoding.EncodeToString(b)
}

func fromArgument(arg string) []byte {
	b, err := base64.RawStdEncoding.DecodeString(arg)
	if err != nil {
		return nil
	}

	return b
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

func EncodeProto(m proto.Message) string {
	return toArgument(msgToByte(m))
}

func DecodeProto(arg string, m proto.Message) error {
	b := fromArgument(arg)
	if b == nil {
		return errors.New("Invalid argument string")
	}

	return proto.Unmarshal(b, m)
}

func (hasher *tx) GenHash(method string) []byte {

	hh := msgToByte(hasher.header)
	hm := msgToByte(hasher.msgObj)
	if hh == nil || hm == nil {
		return nil
	}

	return genHash(hh, hm, method)
}
