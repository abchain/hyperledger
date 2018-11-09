package tx

import (
	"fmt"
	protos "github.com/golang/protobuf/ptypes/empty"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"
)

type InnerTxs CollectiveTxs

func (itxh InnerTxs) TxCall(stub shim.ChaincodeStubInterface,
	function string, args [][]byte) ([]byte, error) {

	function = strings.TrimPrefix(function, ".")

	h, ok := itxh[function]

	if !ok {
		return nil, fmt.Errorf("Chaincode never accept this function [%s]", function)
	}

	//TODO: a trustable list can be checked here against the CCName
	/*
		Discussion: ccname is signed by the tx generator and user should know what they are signing for,
		A possible attacking is replicated the whole tx msg and this re-playing attack will be prevented
		by nonce tracking
	*/

	if len(args) < 2 {
		return nil, fmt.Errorf("Calling arguments is malformed")
	}

	originalFunc := string(args[1])

	//we drop an empty message to pass the unmarshal an unknown original messages
	parser, err := txutil.ParseTx(new(protos.Empty), originalFunc, args[2:])
	if err != nil {
		return nil, err
	}

	return h.txSubCall(stub, function, args[0], parser)
}

//innerTx handler also provide a chaincode interface to handling inner calling
func (itxh InnerTxs) Invoke(stub shim.ChaincodeStubInterface, function string, args [][]byte, _ bool) ([]byte, error) {
	return itxh.TxCall(stub, function, args)
}
