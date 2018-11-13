package rpc

import (
	"errors"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
)

//build a chaincode with single chaincodetx
type dummyCC struct {
	*tx.ChaincodeTx
}

func (c dummyCC) Invoke(stub shim.ChaincodeStubInterface, function string, args [][]byte, _ bool) ([]byte, error) {
	return c.TxCall(stub, function, args)
}

type DummyCallerBuilder struct {
	CCName string
	dummyCC
	caller *ChaincodeAdapter
}

func (d *DummyCallerBuilder) Reset() {
	d.caller = NewLocalChaincode(d)
}

func (d *DummyCallerBuilder) NewTxID(txid string) {

	if d.caller == nil {
		panic("Not init yet")
	}
	d.caller.SpecifyTxID(txid)
}

func (d *DummyCallerBuilder) GetCaller(txid string, h tx.TxHandler) Caller {

	d.dummyCC.ChaincodeTx = &tx.ChaincodeTx{d.CCName, h, nil, nil}
	d.caller.SpecifyTxID(txid)
	return d.caller
}

func (d *DummyCallerBuilder) GetQueryer(h tx.TxHandler) Caller {

	d.dummyCC.ChaincodeTx = &tx.ChaincodeTx{d.CCName, h, nil, nil}
	d.caller.TxIDGen = nil

	return d.caller
}

func (d *DummyCallerBuilder) AppendPreHandler(h tx.TxPreHandler) error {

	if d.caller == nil {
		return errors.New("Not inited")
	}
	d.dummyCC.PreHandlers = append(d.dummyCC.PreHandlers, h)

	return nil
}

func (d *DummyCallerBuilder) Stub() *shim.MockStub {
	if d.caller == nil {
		panic("Builder is not init")
	}

	return d.caller.MockStub
}
