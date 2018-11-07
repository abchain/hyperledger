package fabric_impl

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	public_shim "hyperledger.abchain.org/chaincode/shim"
	"time"
)

type stubAdapter struct {
	shim.ChaincodeStubInterface
}

func (s stubAdapter) GetTxTime() (time.Time, error) {
	ts, err := s.GetTxTimestamp()
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(ts.Seconds, int64(ts.Nanos)), nil
}

func (s stubAdapter) GetRawStub() interface{} {
	return s.ChaincodeStubInterface
}

func (s stubAdapter) RangeQueryState(startKey, endKey string) (public_shim.StateRangeQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.RangeQueryState(startKey, endKey)
}

//for CallerAttr if, function name is different ...
func (s stubAdapter) GetCallerAttribute(attributeName string) ([]byte, error) {
	return s.ReadCertAttribute(attributeName)
}

//for inner invoking
func (s stubAdapter) InvokeChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error) {
	return s.ChaincodeStubInterface.InvokeChaincode(chaincodeName, append([][]byte{[]byte(method)}, args...))
}

func (s stubAdapter) QueryChaincode(chaincodeName string, method string, args [][]byte) ([]byte, error) {
	return s.ChaincodeStubInterface.QueryChaincode(chaincodeName, append([][]byte{[]byte(method)}, args...))
}

func CreateStub(stub shim.ChaincodeStubInterface) public_shim.ChaincodeStubInterface {
	return stubAdapter{stub}
}
