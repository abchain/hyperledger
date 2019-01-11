package service

import (
	"github.com/gocraft/web"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	share "hyperledger.abchain.org/chaincode/modules/sharesubscription"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	tx "hyperledger.abchain.org/core/tx"
	"math/big"
	"strconv"
	"strings"
)

const (
	ContrcatAddr = "contractAddr"
)

type Subscription struct {
	*RPCCoreWithAccount
	token token.GeneralCall
	share share.GeneralCall
}

type SubscriptionRouter struct {
	*web.Router
}

func CreatSubscriptionRouter(root RPCAccountRouter, path string) SubscriptionRouter {
	return SubscriptionRouter{
		root.Subrouter(Subscription{}, path),
	}
}

func (r SubscriptionRouter) Init() SubscriptionRouter {

	r.Middleware((*Subscription).InitCaller)
	return r
}

func (r SubscriptionRouter) BuildRoutes() {

	r.Post("/", (*Subscription).NewContract)
	r.Post("/redeem/:"+ContrcatAddr, (*Subscription).Redeem)
	r.Get("/:"+ContrcatAddr, (*Subscription).QueryContract)
}

func (s *Subscription) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	next(rw, req)
}

type contractEntry struct {
	Txid    string `json:"txID"`
	Address string `json:"contract address"`
}

func (s *Subscription) NewContract(rw web.ResponseWriter, req *web.Request) {

	contractStrs := req.PostForm["contract"]
	contract := make(map[string]int32)

	for _, str := range contractStrs {
		ret := strings.Split(str, ":")
		if len(ret) < 2 {
			s.NormalErrorF(rw, -100, "Wrong contract string")
			return
		}

		w, err := strconv.Atoi(ret[1])
		if err != nil {
			s.NormalError(rw, err)
			return
		}

		contract[ret[0]] = int32(w)
	}

	//s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)
	share := share.GeneralCall{s.TxGenerator}

	var addr *tx.Address
	var err error
	if s.ActivePrivk == nil {
		addr, err = tx.NewAddressFromString(req.PostFormValue("initiator"))
		if err != nil {
			s.NormalError(rw, err)
			return
		}
	} else {
		addr, err = tx.NewAddressFromPrivateKey(s.ActivePrivk)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
	}

	conaddr, err := share.New(contract, addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &contractEntry{
		txid,
		tx.NewAddressFromHash(conaddr).ToString(),
	})
}

func (s *Subscription) Redeem(rw web.ResponseWriter, req *web.Request) {

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
	}

	conaddr, err := tx.NewAddressFromString(req.PathParams[ContrcatAddr])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	share := share.GeneralCall{s.TxGenerator}
	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok && (!amount.IsUint64() || amount.Uint64() != 0) {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	// redeemAddr, err := tx.NewAddress(s.ActivePrivk.Public())
	// if err != nil {
	// 	s.NormalError(rw, err)
	// 	return
	// }

	s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)
	var toAddr *tx.Address
	if tos := req.PostFormValue("to"); tos != "" {
		toAddr, err = tx.NewAddressFromString(tos)
	} else {
		toAddr, err = tx.NewAddressFromPrivateKey(s.ActivePrivk)
	}

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	//we can omit redeem addr in calling message
	nonceid, err := share.Redeem(conaddr.Hash, amount, [][]byte{toAddr.Hash})

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	var nonce []byte
	if len(nonceid.Nonces) > 0 {
		nonce = nonceid.Nonces[0]
	}

	s.Normal(rw, &FundEntry{
		txid,
		s.EncodeEntry(nonce),
		s.TxGenerator.GetBuilder().GetNonce(),
	})

}

type contractMemberEntry struct {
	Weight     float64 `json:"weight"`
	TotalAsset string  `json:"shares"`
	Rest       string  `json:"availiable"`
}

type contractQueryEntry struct {
	Balance    string                          `json:"balance"`
	TotalAsset string                          `json:"shares"`
	Members    map[string]*contractMemberEntry `json:"contract"`
}

func toContractEntry(contract *pb.Contract_s, balance *big.Int) (*contractQueryEntry, error) {

	out := &contractQueryEntry{}
	totalShare := big.NewInt(0).Add(contract.TotalRedeem, balance)

	status := make(map[string]*contractMemberEntry)

	wb := big.NewInt(int64(share.WeightBase))
	for _, s := range contract.Status {

		ret := &contractMemberEntry{}

		canRedeem := big.NewInt(int64(s.Weight))
		canRedeem = canRedeem.Mul(canRedeem, totalShare).Div(canRedeem, wb)

		ret.Weight = float64(s.Weight) / float64(share.WeightBase)
		ret.TotalAsset = canRedeem.String()
		ret.Rest = big.NewInt(0).Sub(canRedeem, s.TotalRedeem).String()

		status[s.MemberID] = ret
	}

	out.Members = status
	out.Balance = balance.String()
	out.TotalAsset = totalShare.String()

	return out, nil
}

func (s *Subscription) QueryContract(rw web.ResponseWriter, req *web.Request) {

	addr, err := tx.NewAddressFromString(req.PathParams[ContrcatAddr])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	tk := token.GeneralCall{s.TxGenerator}
	err, tokenacc := tk.Account(addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	share := share.GeneralCall{s.TxGenerator}
	err, contract := share.Query(addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	ret, err := toContractEntry(contract, tokenacc.Balance)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, ret)
}
