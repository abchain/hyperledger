package service

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	share "hyperledger.abchain.org/chaincode/sharesubscription"
	pb "hyperledger.abchain.org/chaincode/sharesubscription/protos"
	tx "hyperledger.abchain.org/tx"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

const (
	contrcatAddr    = "contractAddr"
	contrcatAddrDir = "/:contractAddr"
)

type subscription struct {
	*apiCore
	token token.GeneralCall
	share share.GeneralCall
}

func (s *subscription) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.token = token.GeneralCall{txgen.SimpleTxGen(chaincode.CC_NAME)}
	s.share = share.GeneralCall{txgen.SimpleTxGen(chaincode.CC_NAME)}
	if offlineMode {
		s.token.Dispatcher = ccCaller
		s.share.Dispatcher = ccCaller
	} else {
		http.Error(rw, "Not implied", http.StatusBadRequest)
		return
	}

	next(rw, req)
}

type contractEntry struct {
	Txid    string `json:"txID"`
	Address string `json:"contract address"`
}

func (s *subscription) NewContract(rw web.ResponseWriter, req *web.Request) {

	if s.activePrivk == nil {
		s.normalErrorF(rw, -100, "No account is specified")
		return
	}

	contractStrs := req.PostForm["contract"]
	contract := make(map[string]uint32)

	for _, str := range contractStrs {
		ret := strings.Split(str, ":")
		if len(ret) < 2 {
			s.normalErrorF(rw, -100, "Wrong contract string")
			return
		}

		w, err := strconv.Atoi(ret[1])
		if err != nil {
			s.normalError(rw, err)
			return
		}

		contract[ret[0]] = uint32(w)
	}

	s.share.Credgenerator = txgen.NewSingleKeyCred(s.activePrivk)

	conaddr, err := s.share.New(contract, s.activePrivk.Public())
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &contractEntry{
		string(s.share.Dispatcher.LastInvokeTxId()),
		tx.NewAddressFromHash(conaddr).ToString(),
	})
}

func (s *subscription) Redeem(rw web.ResponseWriter, req *web.Request) {

	if s.activePrivk == nil {
		s.normalErrorF(rw, -100, "No account is specified")
		return
	}

	conaddr, err := tx.NewAddressFromString(req.PathParams[contrcatAddr])
	if err != nil {
		s.normalError(rw, err)
		return
	}

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok || (amount.IsUint64() && amount.Uint64() == 0) {
		s.normalErrorF(rw, 0, "Invalid amount")
		return
	}

	redeemAddr, err := tx.NewAddress(s.activePrivk.Public())
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.share.Credgenerator = txgen.NewSingleKeyCred(s.activePrivk)

	nonceid, err := s.share.Redeem(conaddr.Hash, redeemAddr.Hash, amount)

	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &fundEntry{
		string(s.share.Dispatcher.LastInvokeTxId()),
		encodeEntry(nonceid),
	})

}

type contractMemberEntry struct {
	Weight     uint32 `json:"weight"`
	TotalAsset string `json:"shares"`
	Rest       string `json:"availiable"`
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
	for addr, s := range contract.Status {

		ret := &contractMemberEntry{}

		haveRedeem := big.NewInt(0).SetBytes(s.TotalRedeem)
		canRedeem := big.NewInt(int64(s.Weight))
		canRedeem = canRedeem.Mul(totalShare, canRedeem).Div(canRedeem, wb)

		ret.Weight = s.Weight
		ret.TotalAsset = canRedeem.String()
		ret.Rest = haveRedeem.Sub(canRedeem, haveRedeem).String()

		status[addr] = ret
	}

	out.Members = status
	out.Balance = addBal.String()
	out.TotalAsset = totalShare.String()

	return out, nil
}

func (s *subscription) QueryContract(rw web.ResponseWriter, req *web.Request) {

	addr, err := tx.NewAddressFromString(req.PathParams[contrcatAddr])
	if err != nil {
		s.normalError(rw, err)
		return
	}

	err, tokenacc := s.token.Account(addr.Hash)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	err, contract := s.share.Query(addr.Hash)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	ret, err := toContractEntry(contract, tokenacc.Balance)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, ret)
}
