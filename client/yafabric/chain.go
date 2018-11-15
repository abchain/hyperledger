package client

import (
	"fmt"
	protos "github.com/abchain/fabric/protos"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/client"
)

type chainAcquire interface {
	GetBlock(int64) *protos.Block
	GetTransaction(string) *protos.Transaction
	GetTxIndex(string) int64
}

type blockchainInterpreter struct {
	chainAcquire
}

func decodeTransactionToInvoke(payload []byte) (*protos.ChaincodeInvocationSpec, error) {

	invoke := &protos.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(payload, invoke); err != nil {
		return nil, fmt.Errorf("protobuf decode fail %s", err.Error())
	}

	if len(invoke.GetChaincodeSpec().GetCtorMsg().GetArgs()) == 0 {
		return nil, fmt.Errorf("Uninitialized invoke tx")
	}

	return invoke, nil
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
	ret.Chaincode = string(tx.GetChaincodeID())
	ret.CreatedFlag = tx.GetType() == protos.Transaction_CHAINCODE_DEPLOY
	inv, err := decodeTransactionToInvoke(tx.GetPayload())
	if err != nil {
		ret.Detail = fmt.Sprintf("Invalid payload (%s)", err)
		return ret
	}

	args := inv.ChaincodeSpec.CtorMsg.Args
	ret.Method = string(args[0])
	ret.Args = args[1:]

	return ret
}

func (i *blockchainInterpreter) GetBlock(h int64) *client.ChainBlock {
	blk := i.chainAcquire.GetBlock(h)
	if blk == nil {
		return nil
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

	return outblk
}

func (i *blockchainInterpreter) GetTransaction(txid string) *client.ChainTransaction {
	tx := i.chainAcquire.GetTransaction(txid)
	if tx == nil {
		return nil
	}
	ret := i.resolveTx(tx)
	ret.Height = i.chainAcquire.GetTxIndex(txid)
	return ret
}

func (i *blockchainInterpreter) GetTxEvent(txid string) *client.ChainTxEvents {
	//no implement
	return nil
}
