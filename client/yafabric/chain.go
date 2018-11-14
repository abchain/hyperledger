package client

import (
	"fmt"
	protos "github.com/abchain/fabric/protos"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/tx"
	"io/ioutil"
	"math/big"
	"strings"
)

type chainAcquire interface {
	GetBlock(int64) *protos.Block
	GetTransaction(string) *protos.Transaction
	GetTxIndex(string) int64
}

type blockchainInterpreter struct {
	chainAcquire
	regParser map[string]client.TxArgParser
}

var notHyperledgerTx = `Not a hyperledger project compatible transaction`

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
	parser, err := abchainTx.ParseTx(new(struct{}), ret.Method, args[1:])
	if err != nil {
		ret.Detail = fmt.Sprintf("Invalid arguments (%s)", err)
		return ret
	}
	ret.Nonce = parser.GetNounce()

	if addParser, ok := i.regParser[strings.Join([]string{ret.Method, ret.Chaincode}, "@")]; ok {
		//a hack: the message is always in args[2]
		msg := addParser.Msg()
		err = proto.Unmarshal(args[2], msg)
		if err != nil {
			ret.Detail = fmt.Sprintf("Invalid message arguments (%s)", err)
			return ret
		}
		ret.Data = msg
		ret.Detail = addParser.Detail(msg)
	}
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

	return outblk
}

func (i *blockchainInterpreter) GetTransaction(tx string) *client.ChainTransaction {
	tx := i.chainAcquire.GetTransaction(tx)
	if tx == nil {
		return nil
	}
	ret := i.resolveTx(tx)
	ret.Height = i.chainAcquire.GetTxIndex(tx)
	return ret
}

func (i *blockchainInterpreter) GetTxEvent(tx string) *client.ChainTxEvents {
	//no implement
	return nil
}

func (i *blockchainInterpreter) RegParser(cc string, method string, p client.TxArgParser) {
	i.regParser[strings.Join([]string{method, cc}, "@")] = p
}
