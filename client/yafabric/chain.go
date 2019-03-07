package client

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/client/yafabric/protos"
)

type chainAcquire interface {
	GetCurrentBlock() (int64, error)
	GetBlock(int64) (*protos.Block, error)
	GetTransaction(string) (*protos.Transaction, error)
	GetTxIndex(string) (int64, error)
}

type blockchainInterpreter struct {
	chainAcquire
}

func decodeTransactionPayload(deployFlag bool, payload []byte) (*protos.ChaincodeSpec, error) {

	var spec *protos.ChaincodeSpec
	if !deployFlag {
		invoke := &protos.ChaincodeInvocationSpec{}
		if err := proto.Unmarshal(payload, invoke); err != nil {
			return nil, fmt.Errorf("protobuf decode invoke fail %s", err.Error())
		}
		spec = invoke.GetChaincodeSpec()
	} else {
		cds := &protos.ChaincodeDeploymentSpec{}
		if err := proto.Unmarshal(payload, cds); err != nil {
			return nil, fmt.Errorf("protobuf decode cds fail %s", err.Error())
		}
		spec = cds.GetChaincodeSpec()
	}

	if len(spec.GetCtorMsg().GetArgs()) == 0 {
		return nil, fmt.Errorf("Uninitialized invoke tx")
	}

	return spec, nil
}

func (i *blockchainInterpreter) resolveTxEvent(txe *protos.ChaincodeEvent) *client.ChainTxEvents {
	ret := new(client.ChainTxEvents)

	ret.TxID = txe.GetTxID()
	ret.Chaincode = txe.GetChaincodeID()
	ret.Name = txe.GetEventName()
	ret.Payload = txe.GetPayload()

	return ret
}

func (i *blockchainInterpreter) resolveTx(tx *protos.Transaction) *client.ChainTransaction {
	ret := new(client.ChainTransaction)

	ret.TxID = tx.GetTxid()
	ret.CreatedFlag = tx.GetType() == protos.Transaction_CHAINCODE_DEPLOY

	if tx.GetConfidentialityLevel() != protos.ConfidentialityLevel_PUBLIC {
		//Can't not parse non-public transaction
		return ret
	}

	spec, err := decodeTransactionPayload(ret.CreatedFlag, tx.GetPayload())
	if err != nil {
		return ret
	}

	ret.Chaincode = spec.GetChaincodeID().GetName()
	args := spec.CtorMsg.Args
	ret.Method = string(args[0])
	ret.TxArgs = args[1:]

	return ret
}

func (i *blockchainInterpreter) GetChain() (*client.Chain, error) {

	if h, err := i.chainAcquire.GetCurrentBlock(); err != nil {
		return nil, err
	} else {
		return &client.Chain{h}, nil
	}

}

func (i *blockchainInterpreter) GetBlock(h int64) (*client.ChainBlock, error) {
	blk, err := i.chainAcquire.GetBlock(h)
	if err != nil {
		return nil, err
	}

	outblk := new(client.ChainBlock)

	for _, tx := range blk.GetTransactions() {
		ret := i.resolveTx(tx)
		ret.Height = h
		outblk.Transactions = append(outblk.Transactions, ret)
	}

	for _, evt := range blk.GetNonHashData().GetChaincodeEvents() {
		ret := i.resolveTxEvent(evt)
		outblk.TxEvents = append(outblk.TxEvents, ret)
	}

	return outblk, nil
}

func (i *blockchainInterpreter) GetTransaction(txid string) (*client.ChainTransaction, error) {
	tx, err := i.chainAcquire.GetTransaction(txid)
	if err != nil {
		return nil, err
	}
	ret := i.resolveTx(tx)
	//not consider as fatal error even fail
	ret.Height, _ = i.chainAcquire.GetTxIndex(txid)
	return ret, nil
}

func (i *blockchainInterpreter) GetTxEvent(txid string) ([]*client.ChainTxEvents, error) {
	//no implement
	return nil, fmt.Errorf("No implement")
}

func (c *RpcClientConfig) Chain() (client.ChainInfo, error) {
	conn, err := c.conn.obtainConn(c.connManager.Context())
	if conn == nil {
		return nil, err
	}

	builder := &RpcBuilder{
		Conn:        *conn,
		ConnManager: c.connManager,
		TxTimeout:   c.TxTimeout,
	}

	if err := builder.VerifyConn(); err != nil {
		return nil, err
	}

	return &blockchainInterpreter{builder}, nil
}
