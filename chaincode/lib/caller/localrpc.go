package rpc

//this tools help you to build a lite caller which wrapping a handler handily

import (
	"hyperledger.abchain.org/chaincode/shim"
	"time"
)

type ChaincodeAdapter struct {
	*shim.MockStub
	LastInvokeId []byte
}

func NewLocalChaincode(cc shim.Chaincode) Caller {
	return &ChaincodeAdapter{MockStub: shim.NewMockStub("LocalCC", cc)}
}

func (c *ChaincodeAdapter) Deploy(method string, arg []string) error {
	_, err := c.MockInit("regtest_deploy_test", method, arg)
	return err
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
	return c.LastInvokeId
}
