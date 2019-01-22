/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package shim provides APIs for the chaincode to access its state
// variables, transaction context and call other chaincodes.
package shim

import (
	"container/list"
	"errors"
	"strings"

	"github.com/op/go-logging"
	"time"
)

// Logger for the shim package.
var mockLogger = logging.MustGetLogger("mock")

// MockStub is an implementation of ChaincodeStubInterface for unit testing chaincode.
// Use this instead of ChaincodeStub in your chaincode's unit test calls to Init, Query or Invoke.
type MockStub struct {
	// arguments the stub was called with
	args [][]byte

	// A pointer back to the chaincode that will invoke this, set by constructor.
	// If a peer calls this stub, the chaincode will be invoked from here.
	cc Chaincode

	// A nice name that can be used for logging
	Name string

	// State keeps name value pairs
	State map[string][]byte

	// Keys stores the list of mapped values in lexical order
	Keys *list.List

	// registered list of other MockStub chaincodes that can be called from this MockStub
	Invokables map[string]*MockStub

	// stores a transaction uuid while being Invoked / Deployed
	// TODO if a chaincode uses recursion this may need to be a stack of TxIDs or possibly a reference counting map
	TxID string

	// store when is called from another chaincode
	InvokedCCName string

	// Extended plugin
	EventHandler func(string, []byte) error
}

func (stub *MockStub) GetRawStub() interface{} {
	return stub
}

func (stub *MockStub) GetTxID() string {
	return stub.TxID
}

func (stub *MockStub) GetArgs() [][]byte {
	return stub.args
}

func (stub *MockStub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

// Used to indicate to a chaincode that it is part of a transaction.
// This is important when chaincodes invoke each other.
// MockStub doesn't support concurrent transactions at present.
func (stub *MockStub) MockTransactionStart(txid string) {
	stub.TxID = txid
}

type stateEntry struct {
	key   string
	value []byte
}

// End a mocked transaction, clearing the UUID.
func (stub *MockStub) MockTransactionEnd(uuid string, e error) {
	if e == nil {
		//done and clear
		stub.Keys = stub.Keys.Init()
	} else {
		mockLogger.Debug("MockStub", stub.Name, "rerolling tx", stub.TxID)

		for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
			elemValue := elem.Value.(*stateEntry)

			if elemValue.value == nil {
				delete(stub.State, elemValue.key)
			} else {
				stub.State[elemValue.key] = elemValue.value
			}

		}

	}
	stub.TxID = ""
}

// Register a peer chaincode with this MockStub
// invokableChaincodeName is the name or hash of the peer
// otherStub is a MockStub of the peer, already intialised
func (stub *MockStub) MockPeerChaincode(invokableChaincodeName string, otherStub *MockStub) {
	stub.Invokables[invokableChaincodeName] = otherStub
}

// Initialise this chaincode,  also starts and ends a transaction.
func (stub *MockStub) MockInit(uuid string, function string, args [][]byte) ([]byte, error) {
	stub.args = getBytes(function, args)
	stub.MockTransactionStart(uuid)
	bytes, err := stub.cc.Invoke(stub, function, args, true)
	stub.MockTransactionEnd(uuid, err)
	return bytes, err
}

// Invoke this chaincode, also starts and ends a transaction.
func (stub *MockStub) MockInvoke(uuid string, function string, args [][]byte) ([]byte, error) {
	stub.args = getBytes(function, args)
	stub.MockTransactionStart(uuid)
	bytes, err := stub.cc.Invoke(stub, function, args, true)
	stub.MockTransactionEnd(uuid, err)
	return bytes, err
}

// Query this chaincode
func (stub *MockStub) MockQuery(function string, args [][]byte) ([]byte, error) {
	stub.args = getBytes(function, args)
	// no transaction needed for queries
	bytes, err := stub.cc.Invoke(stub, function, args, false)
	return bytes, err
}

// GetState retrieves the value for a given key from the ledger
func (stub *MockStub) GetState(key string) ([]byte, error) {
	value := stub.State[key]
	mockLogger.Debug("MockStub", stub.Name, "Getting", key, value)
	return value, nil
}

// PutState writes the specified `value` and `key` into the ledger.
func (stub *MockStub) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		mockLogger.Error("Cannot PutState without a transactions - call stub.MockTransactionStart()?")
		return errors.New("Cannot PutState without a transactions - call stub.MockTransactionStart()?")
	}

	mockLogger.Debug("MockStub", stub.Name, "Putting", key, value)

	val, ok := stub.State[key]
	updatedEntry := &stateEntry{key: key}
	if ok {
		updatedEntry.value = val
	}

	stub.State[key] = value
	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(*stateEntry)
		comp := strings.Compare(key, elemValue.key)
		//		mockLogger.Debug("MockStub", stub.Name, "Compared", key, elemValue, " and got ", comp)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(updatedEntry, elem)
			mockLogger.Debug("MockStub", stub.Name, "Key", key, " inserted before", elem.Value)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			mockLogger.Debug("MockStub", stub.Name, "Key", key, "already in State")
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(updatedEntry)
				mockLogger.Debug("MockStub", stub.Name, "Key", key, "appended")
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(updatedEntry)
		mockLogger.Debug("MockStub", stub.Name, "Key", key, "is first element in list")
	}

	return nil
}

// DelState removes the specified `key` and its value from the ledger.
func (stub *MockStub) DelState(key string) error {
	mockLogger.Debug("MockStub", stub.Name, "Deleting", key, stub.State[key])
	delete(stub.State, key)

	// for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
	// 	if strings.Compare(key, elem.Value.(string)) == 0 {
	// 		stub.Keys.Remove(elem)
	// 	}
	// }

	return nil
}

func (stub *MockStub) RangeQueryState(startKey, endKey string) (StateRangeQueryIteratorInterface, error) {
	return NewMockStateRangeQueryIterator(stub, startKey, endKey), nil
}

// Invokes a peered chaincode.
// E.g. stub1.InvokeChaincode("stub2Hash", funcArgs)
// Before calling this make sure to create another MockStub stub2, call stub2.MockInit(uuid, func, args)
// and register it with stub1 by calling stub1.MockPeerChaincode("stub2Hash", stub2)
func (stub *MockStub) InvokeChaincode(chaincodeName string, function string, args [][]byte) ([]byte, error) {

	otherStub := stub.Invokables[chaincodeName]
	if otherStub == nil {
		mockLogger.Error("Could not find peer chaincode to invoke", chaincodeName)
		return nil, errors.New("Could not find peer chaincode to invoke")
	}
	otherStub.InvokedCCName = stub.Name
	defer func() {
		otherStub.InvokedCCName = ""
	}()
	mockLogger.Debug("MockStub", stub.Name, "Invoking peer chaincode", otherStub.Name, function, args)
	bytes, err := otherStub.MockInvoke(stub.TxID, function, args)
	mockLogger.Debug("MockStub", stub.Name, "Invoked peer chaincode", otherStub.Name, "got", bytes, err)
	return bytes, err
}

func (stub *MockStub) QueryChaincode(chaincodeName string, function string, args [][]byte) ([]byte, error) {

	otherStub := stub.Invokables[chaincodeName]
	if otherStub == nil {
		mockLogger.Error("Could not find peer chaincode to query", chaincodeName)
		return nil, errors.New("Could not find peer chaincode to query")
	}
	otherStub.InvokedCCName = stub.Name
	defer func() {
		otherStub.InvokedCCName = ""
	}()
	mockLogger.Debug("MockStub", stub.Name, "Querying peer chaincode", otherStub.Name, function, args)
	bytes, err := otherStub.MockQuery(function, args)
	mockLogger.Debug("MockStub", stub.Name, "Queried peer chaincode", otherStub.Name, "got", bytes, err)
	return bytes, err
}

func (stub *MockStub) GetCallingChaincodeName() string {
	return stub.InvokedCCName
}

func (stub *MockStub) GetOriginalChaincodeName() string {
	return stub.Name
}

// Not implemented
func (stub *MockStub) GetCallerAttribute(attributeName string) ([]byte, error) {
	return nil, nil
}

// Not implemented
func (stub *MockStub) GetCallerCertificate() ([]byte, error) {
	return nil, nil
}

// Not implemented
func (stub *MockStub) GetBinding() ([]byte, error) {
	return nil, nil
}

// Not implemented
func (stub *MockStub) GetTxTime() (time.Time, error) {
	return time.Now(), nil
}

func (stub *MockStub) SetEvent(name string, payload []byte) error {
	if stub.TxID != "" && stub.EventHandler != nil {
		return stub.EventHandler(name, payload)
	}
	return nil
}

// Constructor to initialise the internal State map
func NewMockStub(name string, cc Chaincode) *MockStub {
	mockLogger.Debug("MockStub(", name, cc, ")")
	s := new(MockStub)
	s.Name = name
	s.cc = cc
	s.State = make(map[string][]byte)
	s.Invokables = make(map[string]*MockStub)
	s.Keys = list.New()

	return s
}

/*****************************
 Range Query Iterator
*****************************/

type MockStateRangeQueryIterator struct {
	Closed   bool
	Stub     *MockStub
	StartKey string
	EndKey   string
	Current  *list.Element
}

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (iter *MockStateRangeQueryIterator) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		mockLogger.Error("HasNext() but already closed")
		return false
	}

	if iter.Current == nil {
		mockLogger.Error("HasNext() couldn't get Current")
		return false
	}

	if iter.Current.Next() == nil {
		// we've reached the end of the underlying values
		mockLogger.Debug("HasNext() but no next")
		return false
	}

	if iter.EndKey == iter.Current.Value {
		// we've reached the end of the specified range
		mockLogger.Debug("HasNext() at end of specified range")
		return false
	}

	mockLogger.Debug("HasNext() got next")
	return true
}

// Next returns the next key and value in the range query iterator.
func (iter *MockStateRangeQueryIterator) Next() (string, []byte, error) {
	if iter.Closed == true {
		mockLogger.Error("MockStateRangeQueryIterator.Next() called after Close()")
		return "", nil, errors.New("MockStateRangeQueryIterator.Next() called after Close()")
	}

	if iter.HasNext() == false {
		mockLogger.Error("MockStateRangeQueryIterator.Next() called when it does not HaveNext()")
		return "", nil, errors.New("MockStateRangeQueryIterator.Next() called when it does not HaveNext()")
	}

	iter.Current = iter.Current.Next()

	if iter.Current == nil {
		mockLogger.Error("MockStateRangeQueryIterator.Next() went past end of range")
		return "", nil, errors.New("MockStateRangeQueryIterator.Next() went past end of range")
	}
	key := iter.Current.Value.(string)
	value, err := iter.Stub.GetState(key)
	return key, value, err
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *MockStateRangeQueryIterator) Close() error {
	if iter.Closed == true {
		mockLogger.Error("MockStateRangeQueryIterator.Close() called after Close()")
		return errors.New("MockStateRangeQueryIterator.Close() called after Close()")
	}

	iter.Closed = true
	return nil
}

func (iter *MockStateRangeQueryIterator) Print() {
	mockLogger.Debug("MockStateRangeQueryIterator {")
	mockLogger.Debug("Closed?", iter.Closed)
	mockLogger.Debug("Stub", iter.Stub)
	mockLogger.Debug("StartKey", iter.StartKey)
	mockLogger.Debug("EndKey", iter.EndKey)
	mockLogger.Debug("Current", iter.Current)
	mockLogger.Debug("HasNext?", iter.HasNext())
	mockLogger.Debug("}")
}

func NewMockStateRangeQueryIterator(stub *MockStub, startKey string, endKey string) *MockStateRangeQueryIterator {
	mockLogger.Debug("NewMockStateRangeQueryIterator(", stub, startKey, endKey, ")")
	iter := new(MockStateRangeQueryIterator)
	iter.Closed = false
	iter.Stub = stub
	iter.StartKey = startKey
	iter.EndKey = endKey
	iter.Current = stub.Keys.Front()

	iter.Print()

	return iter
}

func getBytes(function string, args [][]byte) [][]byte {
	return append([][]byte{[]byte(function)}, args...)

}
