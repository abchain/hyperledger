package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

type ChainTransaction struct {
	Height                         int64 `json:",string"`
	TxID, Chaincode, Method, Nonce string
	CreatedFlag                    bool
	//Data for the original protobuf input (Message part) and Detail left for parser
	Detail, Data interface{} `json:",omitempty"`
}

const (
	TxStatus_Success = 0
	TxStatus_Fail    = 1
)

type ChainTxEvents struct {
	TxID, Chaincode, Name string
	Status                int
	Detail                interface{} `json:",omitempty"`
}

type ChainBlock struct {
	Height       int64 `json:",string"`
	Hash         string
	TimeStamp    string
	Transactions []*ChainTransaction
	TxEvents     []*ChainTxEvents
}

//a parser which can handle the arguments of a transaction with purposed format in hyperledger project
type TxArgParser interface {
	Msg() proto.Message
	Detail(proto.Message) interface{}
}

type ChainClient interface {
	GetBlock(int64) *ChainBlock
	GetTransaction(string) *ChainTransaction
	GetTxEvent(string) *ChainTxEvents
	//TODO, add more methods like get range tx and filter ...

	//Registry parser
	RegParser(string, string, TxArgParser)
}

var ChainProxyViaRPC_Impls map[string]func(RpcClient) ChainClient
var ChainProxy_Impls map[string]func(*viper.Viper) ChainClient
