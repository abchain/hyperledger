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

func defTxidGen() string { return time.Now().String() }

func NewLocalChaincode(cc shim.Chaincode) *ChaincodeAdapter {
	return &ChaincodeAdapter{
		MockStub: shim.NewMockStub("LocalCC", cc),
		TxIDGen:  defTxidGen,
	}
}

func (c *ChaincodeAdapter) SpecifyTxID(txid string) {
	c.TxIDGen = func() string { return txid }
}

func (c *ChaincodeAdapter) DefaultTxID() {
	c.TxIDGen = defTxidGen
}

func (c *ChaincodeAdapter) Deploy(method string, arg [][]byte) (string, error) {

	txid := c.TxIDGen()
	_, err := c.MockInit(txid, method, arg)
	return txid, err
}

func (c *ChaincodeAdapter) Invoke(method string, arg [][]byte) (string, error) {

	txid := c.TxIDGen()
	_, err := c.MockInvoke(txid, method, arg)
	return txid, err
}

func (c *ChaincodeAdapter) Query(method string, arg [][]byte) ([]byte, error) {

	return c.MockQuery(method, arg)
}
