package service

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	mtoken "hyperledger.abchain.org/chaincode/modules/generaltoken/multitoken"
	tokenNonce "hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	tx "hyperledger.abchain.org/core/tx"
	"math/big"
	"strings"
)

const (
	TokenName     = "tokenName"
	TokenNamePath = ":" + TokenName + ":token.\\w+"
	FundID        = "fundID"
	AddressFlag   = "address"
	Simple        = "simple"
)

type Fund struct {
	*util.FabricRPCCore
	token token.TokenTx
}

type FundRouter struct {
	*web.Router
}

func CreateFundRouter(root util.TxRouter, path string) FundRouter {
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
	r.Get("/:"+FundID, (*Fund).QueryTransfer)
}

func (r FundRouter) BuildAddressRoutes() {

	r.Get("/:"+AddressFlag, (*Fund).QueryAddress)
	r.Get("/:"+AccountID+"/:"+AccountIndex, (*Fund).Query)
}

func (r FundRouter) BuildGlobalRoutes() {

	r.Post("/init", (*Fund).InitGlobal)
	r.Post("/", (*Fund).Assign)
	r.Get("/", (*Fund).QueryGlobal)
}

func (s *Fund) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	var err error
	if tname := req.PathParams[TokenName]; tname != "" {
		tname = strings.TrimPrefix(tname, "token.")
		mtoken := &mtoken.GeneralCall{s.TxGenerator}
		s.token, err = mtoken.GetToken(tname)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
	} else {
		s.token = &token.GeneralCall{s.TxGenerator}
	}

	next(rw, req)
}

type FundEntry struct {
	Txid  string `json:"txID,omitempty"`
	Entry string `json:"FundNonce"`
	Nonce []byte `json:"Nonce,omitempty"`
}

func (s *Fund) Fund(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create fund request")
	//token := token.GeneralCall{s.TxGenerator}

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

	var fromAddr *tx.Address
	if s.ActivePrivk == nil {
		fromAddr, err = tx.NewAddressFromString(req.PostFormValue("from"))
	} else {
		fromAddr, err = tx.NewAddressFromPrivateKey(s.ActivePrivk)
	}
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	//s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	nonceid, err := s.token.Transfer(fromAddr.Hash, toAddr.Hash, amount)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		string(txid),
		s.EncodeEntry(nonceid),
		s.TxGenerator.GetBuilder().GetNonce(),
	})
}

func (s *Fund) InitGlobal(rw web.ResponseWriter, req *web.Request) {

	//token deployment
	total, ok := big.NewInt(0).SetString(req.PostFormValue("total"), 0)
	if !ok || total.Int64() == 0 {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	err := s.token.Init(total)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		txid,
		"",
		s.TxGenerator.GetBuilder().GetNonce(),
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

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		string(txid),
		s.EncodeEntry(nonceid),
		s.TxGenerator.GetBuilder().GetNonce(),
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
		data.TotalTokens.String(),
		data.UnassignedTokens.String(),
	})

}

type balanceEntry struct {
	Balance  string `json:"balance"`
	LastFund string `json:"lastFundID"`
}

func (s *Fund) Query(rw web.ResponseWriter, req *web.Request) {

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
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
		data.Balance.String(),
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
		data.Balance.String(),
		s.EncodeEntry(data.LastFund.Noncekey),
	})
}

type fundRecordEntry struct {
	Txid   string `json:"txID"`
	Amount string `json:"amount"`
	fundRecordDetailEntry
}

type fundRecordDetailEntry struct {
	FromLast *FuncRecord `json:"from,omitempty"`
	ToLast   *FuncRecord `json:"to,omitempty"`
	TxTime   string      `json:"txTime,omitempty"`
}

type FuncRecord struct {
	Noncekey string `json:"noncekey"`
	IsSend   bool   `json:"isSend"`
}

func (s *Fund) QueryTransfer(rw web.ResponseWriter, req *web.Request) {

	nc := tokenNonce.GeneralCall{s.TxGenerator}

	nonce, err := s.DecodeEntry(req.PathParams[FundID])

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := nc.Nonce([]byte(nonce))
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	ret := &fundRecordEntry{
		Txid:   data.Txid,
		Amount: data.Amount.String(),
	}

	if req.FormValue(Simple) != "true" {
		ret.FromLast = &FuncRecord{
			s.EncodeEntry(data.FromLast.Noncekey),
			data.FromLast.IsSend,
		}
		ret.ToLast = &FuncRecord{
			s.EncodeEntry(data.ToLast.Noncekey),
			data.ToLast.IsSend,
		}
		ret.TxTime = data.NonceTime.Format("2006-01-02 15:04:05")
	}
	s.Normal(rw, ret)

}
