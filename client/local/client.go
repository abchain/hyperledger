package client

import (
	"fmt"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/client"
	"sync"
	"time"
	"math/rand"
)

type LocalChain struct {
	sync.Mutex
	cc map[string]*rpc.ChaincodeAdapter
	txIndex       map[string]*client.ChainTransaction
	evtIndex      map[string][]*client.ChainTxEvents
	blocks        []*client.ChainBlock
	pendingTxs    []*client.ChainTransaction
	pendingEvents []*client.ChainTxEvents
}

type localCC struct {
	*LocalChain
	ccName  string
	adapter *rpc.ChaincodeAdapter
}

func (c *localCC) innerInvoke(method string, arg [][]byte, createFlag bool) (ret string, e error) {

	c.Lock()
	defer c.Unlock()

	tx := new(client.ChainTransaction)
	tx.TxID = c.adapter.GetTxID()
	tx.Chaincode = c.ccName
	tx.Method = method
	tx.TxArgs = arg
	tx.CreatedFlag = createFlag

	defer func() {
		c.txIndex[tx.TxID] = tx
		if e == nil {
			c.pendingTxs = append(c.pendingTxs, tx)
		}
		c.checkBlock()
	}

	ret, e = c.adapter.Deploy(method, arg)
	return
}

func (c *localCC) Deploy(method string, arg [][]byte) ( string, error) {

	return c.innerInvoke(method, arg, true)
}

func (c *localCC) Invoke(method string, arg [][]byte) ( string,  error) {

	return c.innerInvoke(method, arg, false)

}

func (c *localCC) Query(method string, arg [][]byte) ([]byte, error) {

	c.Lock()
	defer c.Unlock()

	return c.adapter.Query(method, arg)
}

func txidGen() string {
	return ""
}

func (c *LocalChain) AddChaincode(ccName string, cc shim.Chaincode) {
	ccAdapter := rpc.NewLocalChaincode(cc)
	c.setEventHandler(ccName, ccAdapter.MockStub)
	c.cc[ccName] = ccAdapter
}

//now we just make one block - one tx
func (c *LocalChain) checkBlock() {

	if len(c.pendingTxs) == 0 {
		return
	}

	go c.BuildBlock()
}

func (c *LocalChain) BuildBlock() {

	c.Lock()
	defer c.Unlock()	

	blk := new(client.ChainBlock)
	blk.Height = len(c.blocks)
	blk.Hash = "Local"
	blk.TimeStamp = time.Now().String()
	//update indexs
	for _, tx := range c.pendingTxs {
		tx.Height = blk.Height
		c.txIndex[evt.TxID] = tx
	}

	blk.Transactions = c.pendingTxs
	c.pendingTxs = nil

	//also index events 
	for _, evt := range c.pendingEvents {
		tx, ok := c.txIndex[evt.TxID]
		if ok && tx.Height > 0{
			c.evtIndex[evt.TxID] = append(c.evtIndex[evt.TxID], evt)
		}
	}

	blk.TxEvents = c.pendingEvents
	c.pendingEvents = nil

	c.blocks = append(c.blocks, blk)
}

func (c *LocalChain) setEventHandler(ccName string, target *shim.MockStub) error {

	target.EventHandler = func(name string, payload []byte) error {

		eobj := new(client.ChainTxEvents)
		eobj.TxID = target.GetTxID()
		eobj.Chaincode = ccName
		eobj.Name = name
		eobj.Payload = payload

		c.pendingEvents = append(c.pendingEvents, eobj)

		return nil
	}

}

//chaininfo impl
func (c *LocalChain) GetChain() (*client.Chain, error) {

	c.Lock()
	defer c.Unlock()

	return &client.Chain{int64(len(c.blocks))}, nil
}

func (c *LocalChain) GetBlock(i int64) (*client.ChainBlock, error) {

	c.Lock()
	defer c.Unlock()

	if i < 0 || i > len(c.blocks) {
		return nil, fmt.Errorf("Exceed blocknum limit (%d):", len(c.blocks))
	}

	return c.blocks[i], nil
}

func (c *LocalChain) GetTransaction(txid string) (*client.ChainTransaction, error) {
	c.Lock()
	defer c.Unlock()

	return c.txIndex[txid], nil
}

func (c *LocalChain) GetTxEvent(txid string) ([]*client.ChainTxEvents, error) {
	c.Lock()
	defer c.Unlock()

	return c.evtIndex[txid], nil

}
