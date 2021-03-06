package currency

import (
	"bytes"
	"math/big"
	"strings"

	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	mtoken "hyperledger.abchain.org/chaincode/modules/generaltoken/multitoken"
	tokenNonce "hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	tx "hyperledger.abchain.org/core/tx"
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
	//r.Get("/:"+AccountID+"/:"+AccountIndex, (*Fund).Query)
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
	Txid  string      `json:"txID"`
	Nonce string      `json:"Nonce"`
	Entry interface{} `json:"FundNonce,omitempty"`
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
	if addrs := req.PostFormValue("from"); addrs != "" {
		fromAddr, err = tx.NewAddressFromString(addrs)
	} else if s.ActivePrivk != nil {
		fromAddr, err = tx.NewAddressFromPrivateKey(s.ActivePrivk)
	} else {
		s.NormalErrorF(rw, 400, "Unknown source address")
		return
	}
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	if bytes.Compare(fromAddr.Hash, toAddr.Hash) == 0 {
		s.NormalErrorF(rw, 0, "can not transfer asset to the same address")
		return
	}
	//s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	var nonceid interface{}
	if req.PostFormValue("legacy") == "" {
		nonceid, err = s.token.Transfer2(fromAddr.Hash, toAddr.Hash, amount)
	} else {
		nonceid, err = s.token.Transfer(fromAddr.Hash, toAddr.Hash, amount)
	}

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.DefaultOutput(nonceid)
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

	s.DefaultOutput(nil)
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

	s.DefaultOutput(nonceid)

}

func (s *Fund) QueryGlobal(rw web.ResponseWriter, req *web.Request) {

	err, data := s.token.Global()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, data.GetObject())

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

	s.Normal(rw, data.GetObject())
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

	s.Normal(rw, data.GetObject())
}

func (s *Fund) QueryTransfer(rw web.ResponseWriter, req *web.Request) {

	nc := tokenNonce.GeneralCall{s.TxGenerator}

	nonce, err := tokenNonce.NonceKeyFromString(req.PathParams[FundID])

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := nc.Nonce([]byte(nonce))
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, data.GetObject())

}
