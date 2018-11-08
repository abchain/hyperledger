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

func NewRuntime(root string, stub shim.ChaincodeStubInterface, readOnly bool) *ChaincodeRuntime {

	return &ChaincodeRuntime{NewShimMap(root, stub, readOnly), stub, stub}

}
