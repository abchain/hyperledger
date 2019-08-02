package currency

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
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
	*util.FabricRPCCore
}

type SubscriptionRouter struct {
	*web.Router
}

func CreatSubscriptionRouter(root util.TxRouter, path string) SubscriptionRouter {
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

	if s.ActivePrivk != nil {
		s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)
	}

	delegator := req.PostFormValue("delegator")
	if delegator == "" && s.ActivePrivk != nil {
		delegatorAddr, err := tx.NewAddressFromPrivateKey(s.ActivePrivk)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
		delegator = delegatorAddr.ToString()
	}

	share := share.GeneralCall{s.TxGenerator}

	conaddr, err := share.NewByDelegator(contract, delegator)
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

	conaddr := req.PathParams[ContrcatAddr]
	if conaddr == "" {
		s.NormalErrorF(rw, 400, "No contract addr")
		return
	}

	share := share.GeneralCall{s.TxGenerator}
	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok {
		s.NormalErrorF(rw, 0, "Invalid amount")
		return
	}

	if s.ActivePrivk != nil {
		s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)
	}

	var toss []string
	for _, tos := range req.PostForm["to"] {
		toss = append(toss, tos)
	}
	if len(toss) == 0 {
		if s.ActivePrivk == nil {
			s.NormalErrorF(rw, 404, "No redeem address")
			return
		}
		toAddr, err := tx.NewAddressFromPrivateKey(s.ActivePrivk)
		if err != nil {
			s.NormalError(rw, err)
			return
		}
		toss = []string{toAddr.ToString()}
	}

	//we can omit redeem addr in calling message
	nonceids, err := share.Redeem(conaddr, amount, toss)

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.DefaultOutput(nonceids)

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

	wb := big.NewInt(int64(contract.TotalWeight))
	for _, s := range contract.Status {

		ret := &contractMemberEntry{}

		canRedeem := big.NewInt(int64(s.Weight))
		canRedeem = canRedeem.Mul(canRedeem, totalShare).Div(canRedeem, wb)

		ret.Weight = float64(s.Weight) / float64(contract.TotalWeight)
		ret.TotalAsset = canRedeem.String()
		ret.Rest = big.NewInt(0).Sub(canRedeem, s.TotalRedeem).String()

		status[tx.NewAddressFromHash(s.MemberID).ToString()] = ret
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

	if s.ActivePrivk != nil {
		s.TxGenerator.Credgenerator = txgen.NewSingleKeyCred(s.ActivePrivk)
	}

	share := share.GeneralCall{s.TxGenerator}
	var contract *pb.Contract_s
	if isMember := req.FormValue("member"); isMember != "" {

		switch strings.ToLower(isMember) {
		case "n", "no":
			err, contract = share.Query_C(addr.Internal())
		case "y", "yes":
			if s.ActivePrivk == nil {
				s.NormalErrorF(rw, 404, "no query member")
				return
			}
			//try to deduce address from private key
			maddr, err := tx.NewAddress(s.ActivePrivk.Public())
			if err != nil {
				s.NormalError(rw, err)
				return
			}
			isMember = maddr.ToString()
			fallthrough
		default:
			err, contract = share.QueryOne(addr.ToString(), isMember)
		}
	} else {
		err, contract = share.Query_C(addr.Internal())
	}

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
