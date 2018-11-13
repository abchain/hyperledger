package tx

import (
	"fmt"
	protos "github.com/golang/protobuf/ptypes/empty"
	"hyperledger.abchain.org/chaincode/impl"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"
)

type CollectiveTxs_InnerSupport CollectiveTxs

func (itxh CollectiveTxs_InnerSupport) TxCall(stub shim.ChaincodeStubInterface,
	function string, args [][]byte) ([]byte, error) {

	h, ok := itxh[strings.TrimPrefix(function, ".")]

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
func (itxh CollectiveTxs_InnerSupport) Invoke(stub shim.ChaincodeStubInterface, function string, args [][]byte, ro bool) ([]byte, error) {
	if strings.HasPrefix(function, ".") {
		return itxh.TxCall(stub, function, args)
	} else {
		return CollectiveTxs(itxh).Invoke(stub, function, args, ro)
	}

}

type InnerAddrBase struct {
	Root string
	*runtime.Config
}

func (i *InnerAddrBase) rt(stub shim.ChaincodeStubInterface) *runtime.ChaincodeRuntime {
	rrt := runtime.NewRuntime(i.Root, stub, i.Config)
	return rrt.SubRuntime("inneraddr")
}

type InnerAddrRegister struct {
	*InnerAddrBase
	ParseAddress
}

func (v InnerAddrRegister) PreHandling(stub shim.ChaincodeStubInterface, function string, p txutil.Parser) error {
	if !strings.HasPrefix(function, ".") {
		return nil
	}

	ivf, err := impl.GetInnerInvoke(stub)
	if err != nil {
		return err
	}

	rt := v.rt(stub)
	addrs := v.GetAddress().ToString()

	ret, err := rt.Storage.GetRaw(addrs)
	if err != nil {
		return err
	} else if len(ret) > 0 {
		return fmt.Errorf("Registry duplicated address")
	}

	return rt.Storage.SetRaw(addrs, []byte(ivf.GetCallingChaincodeName()))
}

type InnerAddrVerifier struct {
	*InnerAddrBase
	ParseAddress
}

func (v InnerAddrVerifier) PreHandling(stub shim.ChaincodeStubInterface, function string, p txutil.Parser) error {
	if !strings.HasPrefix(function, ".") {
		return nil
	}

	ivf, err := impl.GetInnerInvoke(stub)
	if err != nil {
		return err
	}

	rt := v.rt(stub)

	cc, err := rt.Storage.GetRaw(v.GetAddress().ToString())
	if err != nil {
		return err
	}

	if strings.Compare(string(cc), ivf.GetCallingChaincodeName()) != 0 {
		return fmt.Errorf("Addr is not from registered cc")
	}

	return nil
}
