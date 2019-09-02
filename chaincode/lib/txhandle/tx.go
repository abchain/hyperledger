package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"
)

type TxHandler interface {
	Msg() proto.Message
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
	function string, args [][]byte) ([]byte, error) {

	parser, err := txutil.ParseTx(cci.Handler.Msg(), function, args)
	if err != nil {
		return nil, err
	}

	if strings.Compare(cci.Ccname, parser.GetCCname()) != 0 {
		return nil, errors.New("Unmatched chaincode name")
	}

	//check ts expired
	if parser.GetFlags().IsTimeLock() {
		expT := parser.GetTxTime()
		if expT.IsZero() {
			return nil, errors.New("Enforced timelock Tx has no valid expired time")
		} else if nowT, err := stub.GetTxTime(); err != nil {
			return nil, errors.New("Can not get exec time: " + err.Error())
		} else if nowT.After(expT) {
			return nil, errors.New("Tx is expired")
		}
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

func (cci *ChaincodeTx) txSubCall(stub shim.ChaincodeStubInterface,
	function string, payload []byte, parser txutil.Parser) ([]byte, error) {

	callmsg := cci.Handler.Msg()
	if err := proto.Unmarshal(payload, callmsg); err != nil {
		return nil, err
	}

	originMsg := parser.GetMessage()
	defer parser.UpdateMsg(originMsg)
	parser.UpdateMsg(callmsg)

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
