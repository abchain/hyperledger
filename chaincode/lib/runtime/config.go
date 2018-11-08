package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type DefaultConfig struct {
	RootName string
	ReadOnly bool
}

func (c *DefaultConfig) NewRuntime(stub shim.ChaincodeStubInterface) *ChaincodeRuntime {
	return NewRuntime(c.RootName, stub, c.ReadOnly)
}

func (c *DefaultConfig) NewRuntime_ROflag(stub shim.ChaincodeStubInterface, ro bool) *ChaincodeRuntime {
	return NewRuntime(c.RootName, stub, ro)
}
