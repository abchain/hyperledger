package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/impl"
	"hyperledger.abchain.org/chaincode/shim"
)

type InnerTxGenerator struct {
	impl.InnerInvoke
	txstub        shim.ChaincodeStubInterface
	nonce         []byte
	chaincodeName string
	callError     error
}

type InnerChaincode string

func (s InnerChaincode) NewInnerTxInterface(stub shim.ChaincodeStubInterface, nonce []byte) *InnerTxGenerator {

	v, err := impl.GetInnerInvoke(stub)
	return &InnerTxGenerator{InnerInvoke: v, nonce: nonce, callError: err, chaincodeName: string(s)}
}

func (t *InnerTxGenerator) TxDone() chan struct{} {
	return closedChannel
}

func (t *InnerTxGenerator) Result() TxCallResult {
	return t
}

func (t *InnerTxGenerator) GetNonce() []byte {

	return t.nonce
}

func (t *InnerTxGenerator) TxID() (string, error) { return t.txstub.GetTxID(), t.callError }

//mimic context's implement
var closedChannel = make(chan struct{})

func init() {
	close(closedChannel)
}

//we use sync invoking and return the invoking result, but still allow async query
func (t *InnerTxGenerator) Invoke(method string, msg proto.Message) error {

	if t.InnerInvoke == nil {
		return fmt.Errorf("Invoking not init:", t.callError)
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	callmethod := "." + method

	_, t.callError = t.InnerInvoke.InvokeChaincode(t.chaincodeName, callmethod, append([][]byte{payload}, t.txstub.GetArgs()...))
	return t.callError
}

func (t *InnerTxGenerator) Query(method string, msg proto.Message) (chan QueryResp, error) {

	if t.InnerInvoke == nil {
		return nil, fmt.Errorf("Invoking not init:", t.callError)
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	callmethod := "." + method
	retc := make(chan QueryResp)
	go func() {
		var ret QueryResp
		ret.SuccMsg, ret.ErrMsg = t.InnerInvoke.QueryChaincode(t.chaincodeName, callmethod, append([][]byte{payload}, t.txstub.GetArgs()...))
		retc <- ret
	}()

	return retc, nil
}
