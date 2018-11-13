package impl

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type InnerInvoke interface {

	// InvokeChaincode locally calls the specified chaincode `Invoke` using the
	// same transaction context; that is, chaincode calling chaincode doesn't
	// create a new transaction message.
	InvokeChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error)

	// QueryChaincode locally calls the specified chaincode `Query` using the
	// same transaction context; that is, chaincode calling chaincode doesn't
	// create a new transaction message.
	QueryChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error)

	// acquire the (direct) chaincode name in a inner invoking
	GetCallingChaincodeName() string

	GetOriginalChaincodeName() string
}

var InnerInvokeImpl = []func(stub interface{}) (bool, InnerInvoke){MockStubInvokeImpl}

func GetInnerInvoke(stub shim.ChaincodeStubInterface) (InnerInvoke, error) {
	for _, f := range InnerInvokeImpl {
		if ok, ret := f(stub.GetRawStub()); ok {
			return ret, nil
		}
	}

	return nil, NoImplError
}

func MockStubInvokeImpl(stub interface{}) (bool, InnerInvoke) {
	if s, ok := stub.(*shim.MockStub); ok {
		return true, s
	} else {
		return false, nil
	}
}
