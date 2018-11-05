package rpc

//this tools help you to build a lite caller which wrapping a handler handily

import (
	"hyperledger.abchain.org/chaincode/shim"
	"time"
)

type ChaincodeAdapter struct {
	*shim.MockStub
	TxIDGen func() string
}

func NewLocalChaincode(cc shim.Chaincode) *ChaincodeAdapter {
	return &ChaincodeAdapter{
		MockStub: shim.NewMockStub("LocalCC", cc),
		TxIDGen:  func() string { return time.Now().String() },
	}
}

func (c *ChaincodeAdapter) Deploy(method string, arg []string) (string, error) {
	txid := c.TxIDGen()
	_, err := c.MockInit(txid, method, arg)
	return txid, err
}

func (c *ChaincodeAdapter) Invoke(method string, arg []string) (string, error) {
	txid := c.TxIDGen()
	_, err := c.MockInvoke(txid, method, arg)
	return txid, err
}

func (c *ChaincodeAdapter) Query(method string, arg []string) ([]byte, error) {
	return c.MockQuery(method, arg)
}
