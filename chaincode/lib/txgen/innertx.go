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

//the chaincodename is the name in fabric-level (not the framework's cc name)
type InnerChaincode string

func (s InnerChaincode) NewInnerTxInterface(stub shim.ChaincodeStubInterface, nonce []byte) *InnerTxGenerator {

	v, err := impl.GetInnerInvoke(stub)
	return &InnerTxGenerator{InnerInvoke: v, txstub: stub, nonce: nonce, callError: err, chaincodeName: string(s)}
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

//the "all" args in chaincode is {function, arg ...}
func buildParameters(payload []byte, args [][]byte) ([][]byte, error) {

	if len(args) == 0 {
		//this is malformed but we just pass ...
		return [][]byte{payload}, nil
	}

	fname := string(args[0])
	//detect if we are in another inner invoking ...
	if len(fname) > 0 && fname[0] == '.' {
		if len(args) < 2 {
			return nil, fmt.Errorf("Malformed array of arguments")
		}

		//strip the first arguments (current inner call)
		return append([][]byte{payload}, args[2:]...), nil
	}

	return append([][]byte{payload}, args...), nil
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
	args, err := buildParameters(payload, t.txstub.GetArgs())
	if err != nil {
		return err
	}

	_, t.callError = t.InnerInvoke.InvokeChaincode(t.chaincodeName, callmethod, args)
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
	args, err := buildParameters(payload, t.txstub.GetArgs())
	if err != nil {
		return nil, err
	}
	retc := make(chan QueryResp)
	go func() {
		var ret QueryResp
		ret.SuccMsg, ret.ErrMsg = t.InnerInvoke.QueryChaincode(t.chaincodeName, callmethod, args)
		retc <- ret
	}()

	return retc, nil
}
