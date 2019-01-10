package service

import (
	"math/big"

	"github.com/gocraft/web"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/core/crypto"
	tx "hyperledger.abchain.org/core/tx"
)

const (
	FundID      = "fundID"
	AddressFlag = "address"
	Simple      = "simple"
)

type Fund struct {
	*RPCCoreWithAccount
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

	next(rw, req)
}

type FundEntry struct {
	Txid  string `json:"txID"`
	Entry string `json:"fundNonce"`
	Nonce []byte `json:"Nonce"`
}

func (s *Fund) Fund(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create fund request")
	token := token.GeneralCall{s.TxGenerator}

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

	s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	nonceid, err := token.Transfer(fromAddr.Hash, toAddr.Hash, amount)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := token.Result().TxID()
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

	token := token.GeneralCall{s.TxGenerator}

	//token deployment
	total, ok := big.NewInt(0).SetString(req.PostFormValue("total"), 0)
	if !ok || total.Int64() == 0 {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	err := token.Init(total)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := token.Result().TxID()
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

	token := token.GeneralCall{s.TxGenerator}

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

	nonceid, err := token.Assign(toAddr.Hash, amount)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := token.Result().TxID()
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

	token := token.GeneralCall{s.TxGenerator}

	err, data := token.Global()
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

	token := token.GeneralCall{s.TxGenerator}

	privk, err := s.wallet.LoadPrivKey(req.PathParams[AccountID])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	if indstr, ok := req.PathParams[AccountIndex]; ok {
		index, ok := big.NewInt(0).SetString(indstr, 0)
		if ok {
			privk, err = crypto.GetChildPrivateKey(privk, index)
			if err != nil {
				s.NormalError(rw, err)
				return
			}
		}
	}

	addr, err := tx.NewAddressFromPrivateKey(s.ActivePrivk)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := token.Account(addr.Hash)
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

	token := token.GeneralCall{s.TxGenerator}

	addr, err := tx.NewAddressFromString(req.PathParams[AddressFlag])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := token.Account(addr.Hash)
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
	TxTime   string      `json: "txTime,omitempty"`
}

type FuncRecord struct {
	Noncekey string `json:"noncekey"`
	IsSend   bool   `json:"isSend"`
}

func (s *Fund) QueryTransfer(rw web.ResponseWriter, req *web.Request) {

	token := token.NewFullGeneralCall(s.TxGenerator)

	nonce, err := s.DecodeEntry(req.PathParams[FundID])

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := token.Nonce([]byte(nonce))
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
