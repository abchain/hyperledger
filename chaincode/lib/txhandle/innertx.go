package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	pb "hyperledger.abchain.org/protos"
	"strings"
)

//innerTx handler is a manager which can handling any inner calling
type InnerTxs map[string]*ChaincodeTx

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

	parser, err := txutil.ParseTx(cci.Handler.Msg(), function, args)
	if err != nil {
		return nil, err
	}
}
