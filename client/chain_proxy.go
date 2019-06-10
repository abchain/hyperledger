package client

import (
	"errors"

	"github.com/spf13/viper"
)

type Chain struct {
	Height int64
}

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
	Height       int64  `json:",string"`
	Hash         string `json:",omitempty"`
	PreviousHash string
	TimeStamp    string              `json:",omitempty"`
	Transactions []*ChainTransaction `json:"-"`
	TxEvents     []*ChainTxEvents    `json:"-"`
}

type ChainInfo interface {
	GetChain() (*Chain, error)
	GetBlock(int64) (*ChainBlock, error)
	GetTransaction(string) (*ChainTransaction, error)
	GetTxEvent(string) ([]*ChainTxEvents, error)
}

type ChainClient interface {
	ViaWeb(*viper.Viper) ChainInfo
}

var ChainProxy_Impls = map[string]func() ChainClient{}

func (c *fabricRPCCfg) UseChainREST(name string, vp *viper.Viper) error {

	cli, ok := ChainProxy_Impls[name]
	if !ok {
		return errors.New("No implement")
	}

	c.chain = cli().ViaWeb(vp)

	return nil
}
