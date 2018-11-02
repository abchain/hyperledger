package fabric_impl

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/impl"
)

func FabricStubImpl(stub interface{}) (bool, impl.CallerAttributes) {
	if s, ok := stub.(shim.ChaincodeStubInterface); ok {
		return true, stubAdapter{s}
	} else {
		return false, nil
	}
}

func init() {
	impl.CallerAttrImpl = append(impl.CallerAttrImpl, FabricStubImpl)
}
