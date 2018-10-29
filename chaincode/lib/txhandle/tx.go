package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	txutil "hyperledger.abchain.org/tx"
	"strings"
)

type TxHandler interface {
	Msg() proto.Message
	//	Parse(stub interface{}, method string, args []string) (txutil.Parser, error)
	Call(interface{}, txutil.Parser) ([]byte, error)
}

type TxPreHandler interface {
	PreHandling(interface{}, string, txutil.Parser) error
}

type TxPostHandler interface {
	PostHandling(interface{}, string, txutil.Parser, []byte) ([]byte, error)
}

type ChaincodeTx struct {
	Ccname       string
	Handler      TxHandler
	PreHandlers  []TxPreHandler
	PostHandlers []TxPostHandler
}

func (cci *ChaincodeTx) TxCall(stub interface{},
	function string, args []string) ([]byte, error) {

	parser, err := txutil.ParseTx(cci.Handler.Msg(), function, args)
	if err != nil {
		return nil, err
	}

	if strings.Compare(cci.Ccname, parser.GetCCname()) != 0 {
		return nil, errors.New("Unmatched chaincode name")
	}

	if err != nil {
		return nil, err
	}

	for _, h := range cci.PreHandlers {
		err := h.PreHandling(stub, function, parser)
		if err != nil {
			return nil, err
		}
	}

	ret, err := cci.Handler.Call(stub, parser)

	if err != nil {
		return nil, err
	}

	for _, h := range cci.PostHandlers {
		ret, err = h.PostHandling(stub, function, parser, ret)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
