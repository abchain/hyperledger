package client // import "hyperledger.abchain.org/client"

import (
	"errors"

	"github.com/spf13/viper"
	"time"
)

/*Chain for blockchain */
type Chain struct {
	Height int64
}

/*ChainTransaction for tx */
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

/*ChainTxEvents for events */
type ChainTxEvents struct {
	TxID, Chaincode, Name string
	Status                int
	Payload               []byte `json:"-"`
}

//ChainBlock represent a block concept for fabric-stype blocks
//Notice an transactions which is not accepted (i.e. no contribution
//to world-status) MUST NOT be put in Transactions field
//but one ChainTxEvents corresponding to this failed tx, with Status
//field is non-zero, can be put in TxEvents field
type ChainBlock struct {
	Height       int64  `json:",string"`
	Hash         string `json:",omitempty"`
	PreviousHash string
	TimeStamp    time.Time           `json:",omitempty"`
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
