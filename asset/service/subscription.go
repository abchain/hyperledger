package service

import (
	"github.com/gocraft/web"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	share "hyperledger.abchain.org/chaincode/sharesubscription"
	pb "hyperledger.abchain.org/chaincode/sharesubscription/protos"
	tx "hyperledger.abchain.org/tx"
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

	s.token = token.GeneralCall{s.TxGenerator}
	s.share = share.GeneralCall{TxGenerator: s.TxGenerator}
	s.share.CanOmitRedeemAddr()

	next(rw, req)
}

type contractEntry struct {
	Txid    string `json:"txID"`
	Address string `json:"contract address"`
}

func (s *Subscription) NewContract(rw web.ResponseWriter, req *web.Request) {

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
	}

	contractStrs := req.PostForm["contract"]
	contract := make(map[string]uint32)

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

		contract[ret[0]] = uint32(w)
	}

	s.share.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	conaddr, err := s.share.New(contract, s.ActivePrivk.Public())
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &contractEntry{
		string(s.share.Dispatcher.LastInvokeTxId()),
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

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok && (!amount.IsUint64() || amount.Uint64() != 0) {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	toAddr, err := tx.NewAddressFromString(req.PostFormValue("to"))
	if err != nil {
		toAddr = &tx.Address{}
	}
	// redeemAddr, err := tx.NewAddress(s.ActivePrivk.Public())
	// if err != nil {
	// 	s.NormalError(rw, err)
	// 	return
	// }

	s.share.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)

	//we can omit redeem addr in calling message
	nonceid, err := s.share.Redeem(conaddr.Hash, nil, amount, toAddr.Hash)

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &FundEntry{
		string(s.share.Dispatcher.LastInvokeTxId()),
		s.EncodeEntry(nonceid),
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

func toContractEntry(contract *pb.Contract, balance []byte) (*contractQueryEntry, error) {

	out := &contractQueryEntry{}

	addBal := big.NewInt(0).SetBytes(balance)
	totalShare := big.NewInt(0).SetBytes(contract.TotalRedeem)
	totalShare = totalShare.Add(totalShare, addBal)

	status := make(map[string]*contractMemberEntry)

	wb := big.NewInt(int64(share.WeightBase))
	for _, s := range contract.Status {

		ret := &contractMemberEntry{}

		haveRedeem := big.NewInt(0).SetBytes(s.TotalRedeem)
		canRedeem := big.NewInt(int64(s.Weight))
		canRedeem = canRedeem.Mul(totalShare, canRedeem).Div(canRedeem, wb)

		ret.Weight = float64(s.Weight) / float64(share.WeightBase)
		ret.TotalAsset = canRedeem.String()
		ret.Rest = haveRedeem.Sub(canRedeem, haveRedeem).String()

		status[s.MemberID] = ret
	}

	out.Members = status
	out.Balance = addBal.String()
	out.TotalAsset = totalShare.String()

	return out, nil
}

func (s *Subscription) QueryContract(rw web.ResponseWriter, req *web.Request) {

	addr, err := tx.NewAddressFromString(req.PathParams[ContrcatAddr])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, tokenacc := s.token.Account(addr.Hash)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, contract := s.share.Query(addr.Hash)
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
