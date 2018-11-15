package client

import (
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
)

type ChainTransaction struct {
	Height                  int64 `json:",string"`
	TxID, Chaincode, Method string
	CreatedFlag             bool
	TxArgs                  [][]byte `json:"-"`
}

const (
	TxStatus_Success = 0
	TxStatus_Fail    = 1
)

type ChainTxEvents struct {
	TxID, Chaincode, Name string
	Status                int
	Payload               []byte `json:"-"`
}

type ChainBlock struct {
	Height       int64 `json:",string"`
	Hash         string
	TimeStamp    string
	Transactions []*ChainTransaction
	TxEvents     []*ChainTxEvents
}

type ChainInfo interface {
	GetBlock(int64) *ChainBlock
	GetTransaction(string) *ChainTransaction
	GetTxEvent(string) *ChainTxEvents
}

type ChainClient interface {
	Load(*viper.Viper) (ChainInfo, error)
	ViaRpc(rpc.Caller) (ChainInfo, error)
}

var ChainProxy_Impls map[string]func() ChainClient
