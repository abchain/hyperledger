package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/core/tx"
)

type TxGenerator struct {
	txbuilder     txutil.Builder
	nonce         []byte
	Credgenerator TxCredHandler
	Dispatcher    rpc.Caller
	MethodMapper  map[string]string
	Ccname        string
}

const (
	call_deploy = 0
	call_invoke = 1
	call_query  = 2
)

func (t *TxGenerator) postHandling(method string, callwhich int) ([]byte, error) {

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

	if t.Dispatcher != nil {
		switch callwhich {
		case call_deploy:
			err = t.Dispatcher.Deploy(method, args)
			return nil, err
		case call_invoke:
			return t.Dispatcher.Invoke(method, args)
		case call_query:
			return t.Dispatcher.Query(method, args)
		default:
			panic("Not a calling method")

		}
	}

	return nil, nil
}

func (t *TxGenerator) BeginTx(nonce []byte) {
	t.nonce = nonce
	t.txbuilder = nil
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

func (t *TxGenerator) Invoke(method string, msg proto.Message) ([]byte, error) {
	err := t.txcall(method, msg)
	if err != nil {
		return nil, err
	}

	return t.postHandling(method, call_invoke)
}

func (t *TxGenerator) Query(method string, msg proto.Message) ([]byte, error) {
	err := t.txcall(method, msg)
	if err != nil {
		return nil, err
	}
	return t.postHandling(method, call_query)
}

func (t *TxGenerator) methodName(method string) string {
	m, ok := t.MethodMapper[method]

	if !ok {
		return method
	}

	return m
}

func (t *TxGenerator) GetBuilder() txutil.Builder {
	return t.txbuilder
}

func (t *TxGenerator) MapMethod(m map[string]string) {
	for k, v := range m {
		t.MethodMapper[k] = v
	}
}
