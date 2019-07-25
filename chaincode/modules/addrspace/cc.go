package addrspace

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/shim"
)

type AddressSpace interface {
	RegisterCC() error
	QueryPrefix() ([]byte, error)
	NormalizeAddress([]byte) ([]byte, error)
}

type AddrSpaceConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) AddressSpace
}

type addressSpaceConfig struct {
	Root string
	*runtime.Config
}

const (
	tag_prefix = "AddrSpace_"
)

func NewConfig(tag string) *addressSpaceConfig {
	return &addressSpaceConfig{
		Root:   tag_prefix + tag,
		Config: runtime.NewConfig()}
}

func (gl *addressSpaceConfig) NewTx(stub shim.ChaincodeStubInterface, _ []byte) AddressSpace {

	return &addressSpaceImpl{runtime.NewRuntime(gl.Root, stub, gl.Config)}
}

type dummyImplCfg struct{}

func (dummyImplCfg) NewTx(stub shim.ChaincodeStubInterface, _ []byte) AddressSpace {
	return internalAddrSpaceImpl{}
}

func DummyImplCfg() dummyImplCfg { return dummyImplCfg{} }

func InnerinvokeImpl(cc txgen.InnerChaincode) *innerInvokeConfig {
	return &innerInvokeConfig{InnerChaincode: cc}
}

type innerInvokeConfig struct {
	generalCallCache
	txgen.InnerChaincode
}

func (c *innerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) AddressSpace {
	return GeneralCall{c.NewInnerTxInterface(stub, nc)}.CacheImpl(&c.generalCallCache)
}
