package service

import (
	"encoding/base64"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	tx "hyperledger.abchain.org/tx"
	"math/big"
	"net/http"
)

const (
	fundID         = "fundID"
	fundIDDir      = "/:fundID"
	addressFlag    = "address"
	addressFlagDir = "/:address"
)

type fund struct {
	*apiCore
	token token.GeneralCall
}

func encodeEntry(nonce []byte) string {
	return base64.StdEncoding.EncodeToString(nonce)
}

func decodeEntry(nonce string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(nonce)
}

func (s *fund) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.token = token.GeneralCall{txgen.SimpleTxGen(chaincode.CC_NAME)}
	if offlineMode {
		s.token.Dispatcher = ccCaller
	} else {
		http.Error(rw, "Not implied", http.StatusBadRequest)
		return
	}

	next(rw, req)
}

type fundEntry struct {
	Txid  string `json:"txID"`
	Entry string `json:"fundNonce"`
}

func (s *fund) Fund(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create fund request")

	nonce := req.PostFormValue("nonce")
	if nonce == "" {
		s.token.BeginTx(nil)
	} else {
		s.token.BeginTx([]byte(nonce))
	}

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok || (amount.IsUint64() && amount.Uint64() == 0) {
		s.normalErrorF(rw, 0, "Invalid amount")
		return
	}

	if s.activePrivk == nil {
		s.normalErrorF(rw, -100, "No account is specified")
		return
	}

	fromAddr, err := tx.NewAddressFromPrivateKey(s.activePrivk)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	toAddr, err := tx.NewAddressFromString(req.PostFormValue("to"))
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.token.Credgenerator = txgen.NewSingleKeyCred(s.activePrivk)

	nonceid, err := s.token.Transfer(fromAddr.Hash, toAddr.Hash, amount)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &fundEntry{
		string(s.token.Dispatcher.LastInvokeTxId()),
		encodeEntry(nonceid),
	})
}

func (s *fund) Assign(rw web.ResponseWriter, req *web.Request) {

	amount, ok := big.NewInt(0).SetString(req.PostFormValue("amount"), 0)

	if !ok || (amount.IsUint64() && amount.Uint64() == 0) {
		s.normalErrorF(rw, 0, "Invalid amount")
		return
	}

	toAddr, err := tx.NewAddressFromString(req.PostFormValue("to"))
	if err != nil {
		s.normalError(rw, err)
		return
	}

	nonceid, err := s.token.Assign(toAddr.Hash, amount)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &fundEntry{
		string(s.token.Dispatcher.LastInvokeTxId()),
		encodeEntry(nonceid),
	})

}

type globalEntry struct {
	Total      string `json:"total"`
	Unassigned string `json:"unassign"`
}

func (s *fund) QueryGlobal(rw web.ResponseWriter, req *web.Request) {

	err, data := s.token.Global()
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &globalEntry{
		big.NewInt(0).SetBytes(data.TotalTokens).String(),
		big.NewInt(0).SetBytes(data.UnassignedTokens).String(),
	})

}

type balanceEntry struct {
	Balance  string `json:"balance"`
	LastFund string `json:"lastFundID"`
}

func (s *fund) Query(rw web.ResponseWriter, req *web.Request) {

	privk, err := DefaultWallet.LoadPrivKey(req.PathParams["accountID"])
	if err != nil {
		s.normalError(rw, err)
		return
	}

	index, ok := big.NewInt(0).SetString(req.PathParams["index"], 0)
	if ok {
		privk, err = privk.ChildKey(index)
		if err != nil {
			s.normalError(rw, err)
			return
		}
	}

	addr, err := tx.NewAddressFromPrivateKey(s.activePrivk)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	err, data := s.token.Account(addr.Hash)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &balanceEntry{
		big.NewInt(0).SetBytes(data.Balance).String(),
		encodeEntry(data.LastFund.Noncekey),
	})
}

func (s *fund) QueryAddress(rw web.ResponseWriter, req *web.Request) {

	addr, err := tx.NewAddressFromString(req.PathParams[addressFlag])
	if err != nil {
		s.normalError(rw, err)
		return
	}

	err, data := s.token.Account(addr.Hash)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &balanceEntry{
		big.NewInt(0).SetBytes(data.Balance).String(),
		encodeEntry(data.LastFund.Noncekey),
	})
}

type fundRecordEntry struct {
	Txid   string `json:"txID"`
	Amount string `json:"amount"`
}

func (s *fund) QueryTransfer(rw web.ResponseWriter, req *web.Request) {

	nonce, err := decodeEntry(req.PathParams[fundID])
	if err != nil {
		s.normalError(rw, err)
		return
	}

	err, data := s.token.Nonce([]byte(nonce))
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, &fundRecordEntry{
		data.Txid,
		big.NewInt(0).SetBytes(data.Amount).String(),
	})
}
