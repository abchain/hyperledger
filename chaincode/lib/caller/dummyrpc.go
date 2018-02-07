package rpc

//this tools help you to build a lite caller which wrapping a handler handily

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	"time"
)

type DummyCallerBuilder struct {
	Ccname string
	Stub   shim.ChaincodeStubInterface
}

type dummyCaller struct {
	tx.ChaincodeTx
	stub shim.ChaincodeStubInterface
}

func (d *DummyCallerBuilder) GetCaller(h tx.TxHandler) Caller {

	return &dummyCaller{
		tx.ChaincodeTx{d.Ccname, h, nil, nil},
		d.Stub,
	}
}

func (*DummyCallerBuilder) AppendPreHandler(c Caller, h tx.TxPreHandler) error {

	cc, done := c.(*dummyCaller)

	if !done {
		return errors.New("Not a dummy caller")
	}

	cc.PreHandlers = append(cc.PreHandlers, h)

	return nil
}

func (c *dummyCaller) Invoke(method string, arg []string) ([]byte, error) {
	return c.TxCall(c.stub, method, arg)
}

func (c *dummyCaller) Query(method string, arg []string) ([]byte, error) {
	return c.TxCall(c.stub, method, arg)
}

func (c *dummyCaller) LastInvokeTxId() []byte {
	return []byte(c.stub.GetTxID())
}

type ChaincodeAdapter struct {
	*shim.MockStub
	LastInvokeId []byte
}

func (c *ChaincodeAdapter) Invoke(method string, arg []string) ([]byte, error) {
	txid := time.Now().String()
	c.LastInvokeId = []byte(txid)
	return c.MockInvoke(txid, method, arg)
}

func (c *ChaincodeAdapter) Query(method string, arg []string) ([]byte, error) {
	return c.MockQuery(method, arg)
}

func (c *ChaincodeAdapter) LastInvokeTxId() []byte {
	return []byte(c.LastInvokeId)
}
