package tx

import (
	"github.com/golang/protobuf/proto"
	tx "hyperledger.abchain.org/core/tx"
	"reflect"
)

type genArgParser struct {
	msgGen func() proto.Message
}

type intrinsicDetail interface {
	MsgDetail() interface{}
}

func (g genArgParser) Msg() proto.Message { return g.msgGen() }

func (g genArgParser) Detail(m proto.Message) interface{} {
	if md, ok := m.(intrinsicDetail); ok {
		return md.MsgDetail()
	}
	return m
}

//build tx arg parser from txhandler
func GenerateTxArgParser(m map[string]*ChaincodeTx) map[string]tx.TxArgParser {

	ret := make(map[string]tx.TxArgParser)

	for key, txh := range m {

		testM := txh.Handler.Msg()

		ret[key+"@"+txh.Ccname] = genArgParser{func() proto.Message {
			newM := reflect.New(reflect.ValueOf(testM).Elem().Type())
			return newM.Interface().(proto.Message)
		}}
	}

	return ret

}
