package rpc

import (
	"errors"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
)

//build a chaincode with single chaincodetx
type dummyCC struct {
	tx.ChaincodeTx
}

func (c *dummyCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return c.TxCall(stub, function, args)
}

func (c *dummyCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return c.TxCall(stub, function, args)
}

type DummyCallerBuilder struct {
	CCName string
	cc     *dummyCC
}

func (d *DummyCallerBuilder) GetCaller(h tx.TxHandler) Caller {

	if d.cc != nil {
		d.cc = &dummyCC{tx.ChaincodeTx{d.CCName, h, nil, nil}}
	}

	return NewLocalChaincode(d.cc)
}

func (d *DummyCallerBuilder) AppendPreHandler(h tx.TxPreHandler) error {

	if d.cc == nil {
		return errors.New("Not inited")
	}
	d.cc.PreHandlers = append(d.cc.PreHandlers, h)

	return nil
}
