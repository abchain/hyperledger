package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/core/tx"
	"time"
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
	TxID() (string, error)
}

type TxCaller interface {
	Invoke(method string, msg proto.Message) error
	Query(method string, msg proto.Message) (chan QueryResp, error)
	GetNonce() []byte
	TxDone() chan struct{}
	Result() TxCallResult
}

type TxGenerator struct {
	txbuilder     txutil.Builder
	nonce         []byte
	calledTxid    string
	calledError   error
	callnotify    chan struct{}
	call_deploy   bool
	Timelock      time.Time
	Credgenerator TxCredHandler
	Dispatcher    rpc.Caller
	MethodMapper  map[string]string
	Ccname        string
}

func (t *TxGenerator) postHandling() ([][]byte, error) {

	if t.Dispatcher == nil {
		return nil, errors.New("Dispatcher is not set")
	}

	if t.Credgenerator != nil {
		err := t.Credgenerator.DoCred(t.txbuilder)
		if err != nil {
			return nil, err
		}

		return t.txbuilder.GenArguments()
	} else {
		return t.txbuilder.GenArgumentsWithoutCred()
	}
}

func (t *TxGenerator) txcall(method string, msg proto.Message) error {
	if t.Ccname == "" {
		return errors.New("CC name is not set yet")
	}

	method = t.methodName(method)

	if t.Timelock.IsZero() {

		b, err := txutil.NewTxBuilder(t.Ccname, t.GetNonce(), method, msg)
		if err != nil {
			return err
		}
		t.txbuilder = b

	} else {
		b, err := txutil.NewTxBuilderWithTimeLock(t.Ccname, t.GetNonce(), t.Timelock)
		if err != nil {
			return err
		}

		err = b.SetMessage(msg)
		if err != nil {
			return err
		}
		b.SetMethod(method)
		t.txbuilder = b

	}

	return nil
}

func (t *TxGenerator) methodName(method string) string {
	m, ok := t.MethodMapper[method]

	if !ok {
		return method
	}

	return m
}

func (t *TxGenerator) BeginDeploy(nonce []byte) {
	t.BeginTx(nonce)
	t.SetDeploy()
}

func (t *TxGenerator) SetDeploy() {
	t.call_deploy = true
}

func (t *TxGenerator) BeginTx(nonce []byte) {
	t.nonce = nonce
	t.txbuilder = nil
	t.callnotify = make(chan struct{})
	t.call_deploy = false
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

	return t.callnotify
}

func (t *TxGenerator) Result() TxCallResult {

	<-t.TxDone()

	return t
}

func (t *TxGenerator) GetNonce() []byte {

	if t.nonce == nil {
		t.nonce = txutil.GenerateNonce()
	}

	return t.nonce
}

func (t *TxGenerator) TxID() (string, error) { return t.calledTxid, t.calledError }

func (t *TxGenerator) postinvoke() {
	if t.callnotify != nil {
		close(t.callnotify)
	}
}

func (t *TxGenerator) Invoke(method string, msg proto.Message) error {

	if t.callnotify == nil {
		return errors.New("Must call beginTx before invoking")
	} else {

		select {
		case <-t.callnotify:
			return errors.New("Must call beginTx again before start new invoking")
		default:
		}
	}

	if err := t.txcall(method, msg); err != nil {
		return err
	}

	if args, err := t.postHandling(); err != nil {
		return err
	} else {

		if t.callnotify == nil {
			t.callnotify = closedChannel
		}

		go func(deploy bool) {

			defer t.postinvoke()
			if deploy {
				t.calledTxid, t.calledError = t.Dispatcher.Deploy(method, args)
			} else {
				t.calledTxid, t.calledError = t.Dispatcher.Invoke(method, args)
			}

		}(t.call_deploy)

		return nil
	}

}

func (t *TxGenerator) Query(method string, msg proto.Message) (chan QueryResp, error) {

	err := t.txcall(method, msg)
	if err != nil {
		return nil, err
	}

	if args, err := t.postHandling(); err != nil {
		return nil, err
	} else {

		retc := make(chan QueryResp)

		go func() {
			var ret QueryResp
			ret.SuccMsg, ret.ErrMsg = t.Dispatcher.Query(method, args)
			retc <- ret
		}()

		return retc, nil
	}

}
