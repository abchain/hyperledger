package impl

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type CallerAttributes interface {

	// used to read an specific attribute from the transaction certificate,
	// *attributeName* is passed as input parameter to this function.
	// Example:
	//  attrValue,error:=stub.ReadCertAttribute("position")
	GetCallerAttribute(attributeName string) ([]byte, error)

	// GetCallerCertificate returns caller certificate
	GetCallerCertificate() ([]byte, error)
}

var CallerAttrImpl = []func(stub interface{}) (bool, CallerAttributes){MockStubAttrImpl}

func GetCallerAttributes(stub shim.ChaincodeStubInterface) (CallerAttributes, error) {
	for _, f := range CallerAttrImpl {
		if ok, ret := f(stub.GetRawStub()); ok {
			return ret, nil
		}
	}

	return nil, NoImplError
}

func MockStubAttrImpl(stub interface{}) (bool, CallerAttributes) {
	if s, ok := stub.(*shim.MockStub); ok {
		return true, s
	} else {
		return false, nil
	}
}
