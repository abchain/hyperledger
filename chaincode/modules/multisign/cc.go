package multisign

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	"hyperledger.abchain.org/chaincode/shim"
)

type MultiSignAddressTx interface {
	Contract_C(threshold int32, addrs [][]byte, weights []int32) ([]byte, error)

	// replace the existing multisign addresses specified by 'from' to 'to'
	// if 'to' is empty, corresponding part is just removed, and notice this
	// may cause a contract can not be auth. by anyone again
	Update_C(acc, from, to []byte) error

	Query_C(acc []byte) (error, *pb.Contract_s)
}

type MultiSignConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) MultiSignAddressTx
}

type StandardMultiSignConfig struct {
	Root string
	*runtime.Config
}

func NewConfig(tag string) *StandardMultiSignConfig {
	cfg := runtime.NewConfig()

	return &StandardMultiSignConfig{multisign_tag_prefix + tag, cfg}
}

type baseMultiSignTx struct {
	*runtime.ChaincodeRuntime
	nonce []byte
}

const (
	multisign_tag_prefix = "Multisign_"
)

func (cfg *StandardMultiSignConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) MultiSignAddressTx {
	return &baseMultiSignTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config), nonce}
}

type InnerInvokeConfig struct {
	txgen.InnerChaincode
}

func (c InnerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) MultiSignAddressTx {
	return &GeneralCall{c.NewInnerTxInterface(stub, nc)}
}
