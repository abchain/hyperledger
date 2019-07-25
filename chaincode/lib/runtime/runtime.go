package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
	"time"
)

//simply reorganize the interface in chaincodestub...
type TransactionInterface interface {
	GetArgs() [][]byte
	GetStringArgs() []string
	GetTxID() string
	GetTxTime() (time.Time, error)
}

type CoreInterface interface {
	SetEvent(name string, payload []byte) error
	GetRawStub() interface{}
}

type ChaincodeRuntime struct {
	Storage StateMap
	Tx      TransactionInterface
	Core    CoreInterface
}

func NewRuntime(root string, stub shim.ChaincodeStubInterface, cfg *Config) *ChaincodeRuntime {

	return &ChaincodeRuntime{NewShimMap(root, stub, cfg.ReadOnly), stub, stub}

}

func (r *ChaincodeRuntime) Stub() shim.ChaincodeStubInterface {
	return r.Core.(shim.ChaincodeStubInterface)
}

func (r *ChaincodeRuntime) SubRuntime(node string) *ChaincodeRuntime {
	return &ChaincodeRuntime{r.Storage.SubMap(node), r.Tx, r.Core}

}
