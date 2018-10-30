package shim

//we provide a unify interface for stubs interface provided from different fabric implement
//(ya-fabric, 0.6, 1.x, etc) which is mainly a partial stack from shim interface of 0.6

import (
	"time"
)

// Chaincode interface must be implemented by all chaincodes. The fabric runs
// the transactions by calling these functions as specified.
type Chaincode interface {
	// Init is called during Deploy transaction after the container has been
	// established, allowing the chaincode to initialize its internal data
	Init(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)

	// Invoke is called for every Invoke transactions. The chaincode may change
	// its state variables
	Invoke(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)

	// Query is called for Query transactions. The chaincode may only read
	// (but not modify) its state variables and return the result
	Query(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)
}

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

	// RangeQueryState function can be invoked by a chaincode to query of a range
	// of keys in the state. Assuming the startKey and endKey are in lexical
	// an iterator will be returned that can be used to iterate over all keys
	// between the startKey and endKey, inclusive. The order in which keys are
	// returned by the iterator is random.
	RangeQueryState(startKey, endKey string) (StateRangeQueryIteratorInterface, error)

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

// StateRangeQueryIteratorInterface allows a chaincode to iterate over a range of
// key/value pairs in the state.
type StateRangeQueryIteratorInterface interface {

	// HasNext returns true if the range query iterator contains additional keys
	// and values.
	HasNext() bool

	// Next returns the next key and value in the range query iterator.
	Next() (string, []byte, error)

	// Close closes the range query iterator. This should be called when done
	// reading from the iterator to free up resources.
	Close() error
}
