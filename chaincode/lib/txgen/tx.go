package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/tx"
)

type TxGenerator struct {
	txbuilder     txutil.Builder
	nonce         []byte
	Credgenerator TxCredHandler
	Dispatcher    rpc.Caller
	MethodMapper  map[string]string
	Ccname        string
}

func (t *TxGenerator) postHandling(method string, isInvoke bool) ([]byte, error) {

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
		if isInvoke {
			return t.Dispatcher.Invoke(method, args)
		} else {
			return t.Dispatcher.Query(method, args)
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

	return t.postHandling(method, true)
}

func (t *TxGenerator) Query(method string, msg proto.Message) ([]byte, error) {
	err := t.txcall(method, msg)
	if err != nil {
		return nil, err
	}
	return t.postHandling(method, false)
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
