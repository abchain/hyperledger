package service

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	mtoken "hyperledger.abchain.org/chaincode/modules/generaltoken/multitoken"
	tx "hyperledger.abchain.org/core/tx"
	"math/big"
	"strings"
)

type FundBatch struct {
	*util.TxBatch
}

type FundBatchRouter struct {
	*web.Router
}

func CreateFundBatchRouter(root util.TxBatchRouter, path string) FundBatchRouter {
	return FundBatchRouter{
		root.Subrouter(FundBatch{}, path),
	}
}

func (r FundBatchRouter) BuildRoutes() {

	r.Post("/create", (*FundBatch).InitAndAssign)
}

func (s *FundBatch) capturedTokenName(req *web.Request) string {
	if tname := req.PathParams[TokenName]; tname != "" {
		return strings.TrimPrefix(tname, "token.")
	}

	return ""
}

func (s *FundBatch) InitAndAssign(rw web.ResponseWriter, req *web.Request) {

	var tk token.TokenTx
	var err error

	if n := req.PostFormValue("name"); n != "" {
		mt := &mtoken.GeneralCall{s.AcquireCaller()}
		tk, err = mt.GetToken(n)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
	} else {
		s.NormalErrorF(rw, 0, "No token Name")
	}

	//token deployment
	total, ok := big.NewInt(0).SetString(req.PostFormValue("total"), 0)
	if !ok || total.Int64() == 0 {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	var toAddr *tx.Address
	if s.ActivePrivk == nil {
		toAddr, err = tx.NewAddressFromString(req.PostFormValue("to"))
	} else {
		toAddr, err = tx.NewAddressFromPrivateKey(s.ActivePrivk)
	}
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	if err = tk.Init(total); err != nil {
		s.NormalError(rw, err)
		return
	}

	if nonceid, err := tk.Assign(toAddr.Hash, total); err != nil {
		s.NormalError(rw, err)
		return
	} else {
		s.AddBatchOut(&FundEntry{
			"",
			s.EncodeEntry(nonceid),
			s.TxGenerator.GetNonce(),
		})
	}

}
