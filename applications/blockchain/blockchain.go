package blockchain

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/client"
)

type ChainTransaction struct {
	*client.ChainTransaction
	//Data for the original protobuf input (Message part) and Detail left for parser
	ChaincodeModule, Nonce string
	Detail, Data           interface{} `json:",omitempty"`
}

type ChainTxEvents struct {
	*client.ChainTxEvents
	Detail interface{} `json:",omitempty"`
}

//a parser which can handle the arguments of a transaction with purposed format in hyperledger project
type TxArgParser interface {
	Msg() proto.Message
	Detail(proto.Message) interface{}
}

var registryParsers map[string]TxArgParser
