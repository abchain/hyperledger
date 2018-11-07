package shim

//we provide a unify interface for stubs interface provided from different fabric implement
//(ya-fabric, 0.6, 1.x, etc) which is mainly a partial stack from shim interface of 0.6

import (
	"time"
)

// Chaincode interface purposed to be implemented by all chaincodes. The fabric runs
// the transactions by calling these functions as specified. We use the
// designation mixed with fabric 0.6 and 1.x
// (this interface is not so important like it was in the real fabric implement, just
// to provide a suitable interface in some tools)
type Chaincode interface {
	// Invoke is called for every transactions.
	Invoke(stub ChaincodeStubInterface, function string, args [][]byte, readonly bool) ([]byte, error)
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
