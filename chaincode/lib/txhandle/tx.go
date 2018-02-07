package tx

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	"github.com/golang/protobuf/proto"
	txutil "hyperledger.abchain.org/tx"
	"strings"
)

type TxHandler interface {
	Msg() proto.Message
	//	Parse(stub shim.ChaincodeStubInterface, method string, args []string) (txutil.Parser, error)
	Call(shim.ChaincodeStubInterface, txutil.Parser) ([]byte, error)
}

type TxPreHandler interface {
	PreHandling(shim.ChaincodeStubInterface, string, txutil.Parser) error
}

type TxPostHandler interface {
	PostHandling(shim.ChaincodeStubInterface, string, txutil.Parser, []byte) ([]byte, error)
}

type ChaincodeTx struct {
	Ccname       string
	Handler      TxHandler
	PreHandlers  []TxPreHandler
	PostHandlers []TxPostHandler
}

func (cci *ChaincodeTx) TxCall(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {

	parser, err := txutil.ParseTx(cci.Handler.Msg(), stub, function, args)
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
