package service

import (
	"github.com/gocraft/web"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	tx "hyperledger.abchain.org/tx"
	"math/big"
)

const (
	FundID      = "fundID"
	AddressFlag = "address"
)

type Fund struct {
	*RPCCoreWithAccount
	token token.GeneralCall
}

type FundRouter struct {
	*web.Router
}

func CreateFundRouter(root RPCAccountRouter, path string) FundRouter {
	return FundRouter{
		root.Subrouter(Fund{}, path),
	}
}

func (r FundRouter) Init() FundRouter {

	r.Middleware((*Fund).InitCaller)
	return r
}

func (r FundRouter) BuildFundRoutes() {

	r.Post("/", (*Fund).Fund)
	r.Get(":/"+FundID, (*Fund).QueryTransfer)
}

func (r FundRouter) BuildAddressRoutes() {

	r.Get(":/"+AddressFlag, (*Fund).QueryAddress)
	r.Get(":/"+AccountID+":/"+AccountIndex, (*Fund).Query)
}

func (r FundRouter) BuildGlobalRoutes() {

	r.Post("/", (*Fund).Assign)
	r.Get("/", (*Fund).QueryGlobal)
}

func (s *Fund) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.token = token.GeneralCall{s.TxGenerator}
	next(rw, req)
}

type FundEntry struct {
	Txid  string `json:"txID"`
	Entry string `json:"fundNonce"`
}

func (s *Fund) Fund(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create fund request")

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok || (amount.IsUint64() && amount.Uint64() == 0) {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
	}

	fromAddr, err := tx.NewAddressFromPrivateKey(s.ActivePrivk)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	toAddr, err := tx.NewAddressFromString(req.PostFormValue("to"))
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.token.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	nonceid, err := s.token.Transfer(fromAddr.Hash, toAddr.Hash, amount)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		string(s.token.Dispatcher.LastInvokeTxId()),
		s.EncodeEntry(nonceid),
	})
}

func (s *Fund) Assign(rw web.ResponseWriter, req *web.Request) {

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok || (amount.IsUint64() && amount.Uint64() == 0) {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	toAddr, err := tx.NewAddressFromString(req.PostFormValue("to"))
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	nonceid, err := s.token.Assign(toAddr.Hash, amount)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		string(s.token.Dispatcher.LastInvokeTxId()),
		s.EncodeEntry(nonceid),
	})

}

type globalEntry struct {
	Total      string `json:"total"`
	Unassigned string `json:"unassign"`
}

func (s *Fund) QueryGlobal(rw web.ResponseWriter, req *web.Request) {

	err, data := s.token.Global()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &globalEntry{
		big.NewInt(0).SetBytes(data.TotalTokens).String(),
		big.NewInt(0).SetBytes(data.UnassignedTokens).String(),
	})

}

type balanceEntry struct {
	Balance  string `json:"balance"`
	LastFund string `json:"lastFundID"`
}

func (s *Fund) Query(rw web.ResponseWriter, req *web.Request) {

	privk, err := s.wallet.LoadPrivKey(req.PathParams[AccountID])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	index, ok := big.NewInt(0).SetString(req.PathParams[AccountIndex], 0)
	if ok {
		privk, err = privk.ChildKey(index)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
	}

	addr, err := tx.NewAddressFromPrivateKey(s.ActivePrivk)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := s.token.Account(addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &balanceEntry{
		big.NewInt(0).SetBytes(data.Balance).String(),
		s.EncodeEntry(data.LastFund.Noncekey),
	})
}

func (s *Fund) QueryAddress(rw web.ResponseWriter, req *web.Request) {

	addr, err := tx.NewAddressFromString(req.PathParams[AddressFlag])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := s.token.Account(addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &balanceEntry{
		big.NewInt(0).SetBytes(data.Balance).String(),
		s.EncodeEntry(data.LastFund.Noncekey),
	})
}

type fundRecordEntry struct {
	Txid   string `json:"txID"`
	Amount string `json:"amount"`
}

func (s *Fund) QueryTransfer(rw web.ResponseWriter, req *web.Request) {

	nonce, err := s.DecodeEntry(req.PathParams[FundID])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := s.token.Nonce([]byte(nonce))
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &fundRecordEntry{
		data.Txid,
		big.NewInt(0).SetBytes(data.Amount).String(),
	})
}
