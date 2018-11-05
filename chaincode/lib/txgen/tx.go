package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/core/tx"
)

type QueryResp struct {
	SuccMsg []byte
	ErrMsg  error
}

func SyncQueryResult(msg proto.Message, out chan QueryResp) error {

	resp := <-out
	if resp.ErrMsg != nil {
		return resp.ErrMsg
	} else {
		return rpc.DecodeRPCResult(msg, resp.SuccMsg)
	}
}

type TxCallResult interface {
	Nonce() []byte
	TxID() (string, error)
}

type TxCaller interface {
	Invoke(method string, msg proto.Message) error
	Query(method string, msg proto.Message) (chan QueryResp, error)
	TxDone() chan struct{}
	Result() TxCallResult
}

//mimic context's implement
var closedChannel = make(chan struct{})

func init() {
	close(closedChannel)
}

type TxGenerator struct {
	txbuilder     txutil.Builder
	nonce         []byte
	calledTxid    string
	calledError   error
	callnotify    chan struct{}
	call_method   int
	Credgenerator TxCredHandler
	Dispatcher    rpc.Caller
	MethodMapper  map[string]string
	Ccname        string
}

const (
	call_invoke = 0
	call_deploy = 1
	call_query  = 2
)

func (t *TxGenerator) postHandling(method string) (*QueryResp, error) {

	var args []string
	var err error

	if t.Credgenerator != nil {
		err = t.Credgenerator.DoCred(t.txbuilder)
		if err != nil {
			return nil, err
		}

		args, err = t.txbuilder.GenArguments()
	} else {
		args, err = t.txbuilder.GenArgumentsWithoutCred()
	}

	if err != nil {
		return nil, err
	}

	postinvoke := func() {
		if t.callnotify != nil && t.callnotify != closedChannel {
			close(t.callnotify)
		}
		t.callnotify = closedChannel
	}

	//TODO: curreltly we just use a sync calling ....
	if t.Dispatcher != nil {
		switch t.call_method {
		case call_deploy:
			defer postinvoke()
			t.calledTxid, t.calledError = t.Dispatcher.Deploy(method, args)
			return nil, nil
		case call_invoke:
			defer postinvoke()
			t.calledTxid, t.calledError = t.Dispatcher.Invoke(method, args)
			return nil, nil
		case call_query:
			var ret []byte
			ret, err = t.Dispatcher.Query(method, args)
			return &QueryResp{ret, err}, nil
		default:
			panic("Not a calling method")

		}
	}

	return nil, nil
}

func (t *TxGenerator) txcall(method string, msg proto.Message) error {
	if t.Ccname == "" {
		return errors.New("CC name is not set yet")
	}

	method = t.methodName(method)

	b, err := txutil.NewTxBuilder(t.Ccname, t.nonce, method, msg)
	if err != nil {
		return err
	}

	t.txbuilder = b

	return nil
}

func (t *TxGenerator) methodName(method string) string {
	m, ok := t.MethodMapper[method]

	if !ok {
		return method
	}

	return m
}

func (t *TxGenerator) SetDeploy() {
	t.call_method = call_deploy
}

func (t *TxGenerator) BeginTx(nonce []byte) {
	t.nonce = nonce
	t.txbuilder = nil
	t.callnotify = nil
	t.call_method = call_invoke
}

func (t *TxGenerator) GetBuilder() txutil.Builder {
	return t.txbuilder
}

func (t *TxGenerator) MapMethod(m map[string]string) {
	for k, v := range m {
		t.MethodMapper[k] = v
	}
}

func (t *TxGenerator) TxDone() chan struct{} {

	if t.callnotify == nil {
		t.callnotify = make(chan struct{})
	}

	return t.callnotify
}

func (t *TxGenerator) Result() TxCallResult {
	return t
}

func (t *TxGenerator) Nonce() []byte { return t.txbuilder.GetNonce() }

func (t *TxGenerator) TxID() (string, error) { return t.calledTxid, t.calledError }

func (t *TxGenerator) Invoke(method string, msg proto.Message) error {

	if t.callnotify == closedChannel {
		return errors.New("Must call beginTx before new invoking")
	}

	err := t.txcall(method, msg)
	if err != nil {
		return err
	}

	_, err = t.postHandling(method)
	return err
}

func (t *TxGenerator) Query(method string, msg proto.Message) (chan QueryResp, error) {
	t.call_method = call_query
	err := t.txcall(method, msg)
	if err != nil {
		return nil, err
	}
	ret, err := t.postHandling(method)
	if err != nil {
		return nil, err
	}

	retc := make(chan QueryResp, 1)
	retc <- *ret
	return retc, nil
}
