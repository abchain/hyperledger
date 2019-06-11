package hyfabric_impl

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	public_shim "hyperledger.abchain.org/chaincode/shim"
)

type stubAdapter struct {
	shim.ChaincodeStubInterface
}

type mockStateRangeQueryIterator struct {
	shim.MockQueryIteratorInterface
}

func (s stubAdapter) GetTxTime() (time.Time, error) {
	ts, err := s.GetTxTimestamp()
	if err != nil {
		return time.Unix(0, 0).UTC(), err
	}

	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC(), nil
}

func (s stubAdapter) GetRawStub() interface{} {
	return s.ChaincodeStubInterface
}

func (s stubAdapter) RangeQueryState(startKey, endKey string) (public_shim.StateRangeQueryIteratorInterface, error) {
	srqi, err := s.ChaincodeStubInterface.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	return NewMockStateRangeQueryIterator(srqi), nil
	// return nil, nil
}

//for CallerAttr if, function name is different ...
func (s stubAdapter) GetCallerAttribute(attributeName string) ([]byte, error) {
	return nil, nil
}
func (s stubAdapter) GetCallerCertificate() ([]byte, error) {
	return nil, nil
}

//for inner invoking 缺少channelName
func (s stubAdapter) InvokeChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error) {
	channelID := s.ChaincodeStubInterface.GetChannelID()
	fmt.Println("获取的channelID 是：", channelID)
	pb := s.ChaincodeStubInterface.InvokeChaincode(chaincodeName, append([][]byte{[]byte(method)}, args...), channelID)
	return pb.Payload, nil
	// return s.ChaincodeStubInterface.InvokeChaincode(chaincodeName, append([][]byte{[]byte(method)}, args...))
}

func (s stubAdapter) QueryChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error) {
	return s.InvokeChaincode(chaincodeName, method, args)
}

const notInnerInvoking = "NotInnerInvoking"

func (s stubAdapter) GetCallingChaincodeName() string {
	return notInnerInvoking
}

func (s stubAdapter) GetOriginalChaincodeName() string {
	return notInnerInvoking
}

func CreateStub(stub shim.ChaincodeStubInterface) public_shim.ChaincodeStubInterface {
	return stubAdapter{stub}
}

func NewMockStateRangeQueryIterator(result shim.StateQueryIteratorInterface) public_shim.StateRangeQueryIteratorInterface {
	m := new(mockStateRangeQueryIterator)
	m.MockQueryIteratorInterface = result
	return m
}
func (m *mockStateRangeQueryIterator) HasNext() bool {
	return m.MockQueryIteratorInterface.HasNext()
}

func (m *mockStateRangeQueryIterator) Next() (string, []byte, error) {
	kv, err := m.MockQueryIteratorInterface.Next()
	return kv.Key, kv.Value, err
}

func (m *mockStateRangeQueryIterator) Close() error {
	return m.MockQueryIteratorInterface.Close()
}
