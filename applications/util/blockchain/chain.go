package blockchain

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/client"
	"net/http"
	"strconv"
)

type FabricBlockChain struct {
	*util.FabricRPCBase
	cli client.ChainInfo
}

type BlockChainRouter struct {
	*web.Router
}

func CreateBlockChainRouter(root util.RPCRouter, path string) BlockChainRouter {
	return BlockChainRouter{
		root.Subrouter(FabricBlockChain{}, path),
	}
}

const (
	FabricProxy_BlockHeight   = "height"
	FabricProxy_TransactionID = "txID"
)

func (r BlockChainRouter) Init(cf func() (client.ChainInfo, error)) BlockChainRouter {
	r.Middleware(func(s *FabricBlockChain, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		var err error
		s.cli, err = cf()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		next(rw, req)
	})

	return r
}

func (r BlockChainRouter) BuildRoutes() {
	r.Get("/", (*FabricBlockChain).GetBlockchainInfo)
	r.Get("/blocks/:"+FabricProxy_BlockHeight, (*FabricBlockChain).GetBlock)
	r.Get("/transactions/:"+FabricProxy_TransactionID, (*FabricBlockChain).GetTransaction)
	r.Post("/parseTx", (*FabricBlockChain).ParseTransaction)
}

var notHyperledgerTx = `Not a hyperledger project compatible transaction`
var noParser = `No parser can be found for this transaction/event`

type restError struct {
	Error string `json:"Error,omitempty"`
}

func (s *FabricBlockChain) GetBlockchainInfo(rw web.ResponseWriter, req *web.Request) {

	if c, err := s.cli.GetChain(); err != nil {
		s.NormalError(rw, err)
	} else {
		s.Normal(rw, c)
	}
}

func (s *FabricBlockChain) GetBlock(rw web.ResponseWriter, req *web.Request) {

	h, err := strconv.ParseInt(req.PathParams[FabricProxy_BlockHeight], 10, 64)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	blk, err := s.cli.GetBlock(h)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	ret := &ChainBlock{blk, make([]*ChainTransaction, len(blk.Transactions)), make([]*ChainTxEvents, len(blk.TxEvents))}

	for i, tx := range blk.Transactions {
		ret.Transactions[i] = handleTransaction(tx)
	}

	for i, tve := range blk.TxEvents {
		ret.TxEvents[i] = handleTxEvent(tve)
	}

	s.Normal(rw, ret)

}

func (s *FabricBlockChain) GetTransaction(rw web.ResponseWriter, req *web.Request) {

	transactionID := req.PathParams[FabricProxy_TransactionID]

	if transactionID == "" {
		s.NormalErrorF(rw, -100, "Invalid TxID string")
		return
	}

	tx, err := s.cli.GetTransaction(transactionID)
	if err != nil {
		s.NormalError(rw, err)
		return
	} else if tx == nil {
		s.NormalErrorF(rw, 404, "Transaction not found")
		return
	}

	s.Normal(rw, handleTransaction(tx))
}

func (s *FabricBlockChain) ParseTransaction(rw web.ResponseWriter, req *web.Request) {

	var err error
	var txflag string
	tx := new(client.ChainTransaction)

	err, txflag, tx.Method, tx.Chaincode, tx.TxArgs = util.ParseCompactFormTx(req.PostFormValue("tx"))

	if txflag == "D" {
		tx.CreatedFlag = true
	}

	tx.TxID = "Unknown"
	tx.Chaincode = "Unknown"

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, handleTransaction(tx))
}
