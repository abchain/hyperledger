package shim

//we provide a unify interface for stubs interface provided from different fabric implement
//(ya-fabric, 0.6, 1.x, etc) which is mainly a partial stack from shim interface of 0.6

import (
	"time"
)

type ChaincodeStubInterface interface {
	// Get the arguments to the stub call as a 2D byte array
	GetArgs() [][]byte

	// Get the arguments to the stub call as a string array
	GetStringArgs() []string

	// Get the transaction ID
	GetTxID() string

	// returns transaction created timestamp, which is currently
	// taken from the peer receiving the transaction. Note that this timestamp
	// may not be the same with the other peers' time.
	GetTxTime() (time.Time, error)

	// InvokeChaincode locally calls the specified chaincode `Invoke` using the
	// same transaction context; that is, chaincode calling chaincode doesn't
	// create a new transaction message.
	InvokeChaincode(chaincodeName string, args [][]byte) ([]byte, error)

	// QueryChaincode locally calls the specified chaincode `Query` using the
	// same transaction context; that is, chaincode calling chaincode doesn't
	// create a new transaction message.
	QueryChaincode(chaincodeName string, args [][]byte) ([]byte, error)

	// GetState returns the byte array value specified by the `key`.
	GetState(key string) ([]byte, error)

	// PutState writes the specified `value` and `key` into the ledger.
	PutState(key string, value []byte) error

	// DelState removes the specified `key` and its value from the ledger.
	DelState(key string) error

	// used to read an specific attribute from the transaction certificate,
	// *attributeName* is passed as input parameter to this function.
	// Example:
	//  attrValue,error:=stub.ReadCertAttribute("position")
	GetCallerAttribute(attributeName string) ([]byte, error)

	// GetCallerCertificate returns caller certificate
	GetCallerCertificate() ([]byte, error)

	// GetBinding returns the transaction binding
	GetBinding() ([]byte, error)

	// SetEvent saves the event to be sent when a transaction is made part of a block
	SetEvent(name string, payload []byte) error

	// obtain the original chaincodestub interface for more implement-spec code
	GetRawStub() interface{}
}
