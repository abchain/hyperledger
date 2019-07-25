package tx

import (
	"fmt"
	protos "github.com/golang/protobuf/ptypes/empty"
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

//--------------- DO NOT USE, HAVE BEEN DEPRECATED, SEE ADDRSPACDE MODULE ----------------------

//this module help to bind a specified addr with a chaincode, so other chaincode
//could not touch this address

/*
type InnerAddrBase struct {
	Root string
	*runtime.Config
}

func (i *InnerAddrBase) getRT(stub shim.ChaincodeStubInterface) *runtime.ChaincodeRuntime {
	rrt := runtime.NewRuntime(i.Root, stub, i.Config)
	return rrt.SubRuntime("inneraddr")
}

type InnerAddrRegister struct {
	*InnerAddrBase
	ListAddresses
}

func (v InnerAddrRegister) PostHandling(stub shim.ChaincodeStubInterface, function string, p txutil.Parser, retbt []byte) ([]byte, error) {
	if !strings.HasPrefix(function, ".") {
		return retbt, nil
	}

	ivf, err := impl.GetInnerInvoke(stub)
	if err != nil {
		return nil, err
	}

	var addrs []string
	if v.ListAddresses != nil {
		for _, addr := range v.ListAddresses(p.GetMessage()) {
			addrs = append(addrs, addr.ToString())
		}
	} else if maddr, ok := p.GetMessage().(MsgAddresses); ok {
		for _, addr := range maddr.GetAddresses() {
			addrs = append(addrs, addr.ToString())
		}
	}

	if len(addrs) == 0 {
		//not consider as error
		return retbt, nil
	}

	rt := v.getRT(stub)
	dataToSet := []byte(ivf.GetCallingChaincodeName())

	for _, addr := range addrs {
		ret, err := rt.Storage.GetRaw(addr)
		if err != nil {
			return nil, err
		} else if len(ret) > 0 {
			return nil, fmt.Errorf("Registry duplicated address")
		}

		err = rt.Storage.SetRaw(addr, dataToSet)
		if err != nil {
			return nil, err
		}
	}

	return retbt, nil
}

type InnerAddrVerifier struct {
	*InnerAddrBase
	callingccName string
	rt            *runtime.ChaincodeRuntime
}

func (v *InnerAddrVerifier) Clone() AddrVerifier {
	return &InnerAddrVerifier{InnerAddrBase: v.InnerAddrBase}
}

func (v *InnerAddrVerifier) Verify(addr *txutil.Address) error {
	if v.callingccName == "" {
		return nil
	}

	cc, err := v.rt.Storage.GetRaw(addr.ToString())
	if err != nil {
		return err
	}

	if strings.Compare(string(cc), v.callingccName) != 0 {
		return fmt.Errorf("Addr is not from registered cc")
	}

	return nil
}

func (v *InnerAddrVerifier) PreHandling(stub shim.ChaincodeStubInterface, function string, p txutil.Parser) error {
	if !strings.HasPrefix(function, ".") {
		v.callingccName = ""
		return nil
	}

	ivf, err := impl.GetInnerInvoke(stub)
	if err != nil {
		return err
	}

	v.callingccName = ivf.GetCallingChaincodeName()
	v.rt = v.getRT(stub)

	return nil
}
*/
