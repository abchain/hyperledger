package fabric_impl

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/impl"
)

func FabricStubInvokeImpl(stub interface{}) (bool, impl.InnerInvoke) {
	if s, ok := stub.(shim.ChaincodeStubInterface); ok {
		return true, stubAdapter{s}
	} else {
		return false, nil
	}
}

func init() {
	impl.InnerInvokeImpl = append(impl.InnerInvokeImpl, FabricStubInvokeImpl)
}
