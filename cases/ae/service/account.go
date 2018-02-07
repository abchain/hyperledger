package service

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/crypto"
	"hyperledger.abchain.org/tx"
	"math/big"
	"net/http"
)

const (
	accountID       = "accountID"
	accountIDDir    = "/:accountID"
	accountIndex    = "index"
	accountIndexDir = "/:index"
)

type account struct {
	*apiCore
	accountID     string
	shouldPersist bool
}

func (s *account) PersistAccount(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	next(rw, req)

	if s.shouldPersist {
		DefaultWallet.Persist()
	}
}

func (s *account) ParseParameters(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	//add parseform action for PATCH method
	if req.Method == http.MethodPatch {
		err := req.ParseForm()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//but we only get accountID from post method
	if req.Method == http.MethodPost {
		s.accountID = req.PostFormValue(accountID)
	} else {
		s.accountID = req.PathParams[accountID]
	}

	next(rw, req)
}

func (s *account) Create(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create account request")

	// Check accountID
	if s.accountID == "" {
		s.normalError(rw, errors.New("Must provide accountID"))
		return
	}

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Create private key
	priv, err := DefaultWallet.NewPrivKey(s.accountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	// Debug
	logger.Debugf("output: address(%v)", addr)

	s.normal(rw, addr.ToString())

	// Persist wallet data
	s.shouldPersist = true
}

func (s *account) List(rw web.ResponseWriter, req *web.Request) {
	logger.Debug("Received list accounts request")

	// get account list
	privkeys, err := DefaultWallet.ListAll()
	if err != nil {
		s.normalError(rw, err)
		return
	}

	ret := make(map[string]string)
	for k, v := range privkeys {
		addr, err := abchainTx.NewAddressFromPrivateKey(v)
		if err != nil {
			continue
		}

		ret[k] = addr.ToString()
	}

	// Debug
	logger.Debugf("output: address map(%v)", ret)

	s.normal(rw, ret)
}

func (s *account) Query(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received query account request")

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Get address
	priv, err := DefaultWallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.normalErrorF(rw, 404, "account Not Found")
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, addr.ToString())
}

func (s *account) QueryChild(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received query child account request")

	// Parse index
	index, ok := big.NewInt(0).SetString(req.PathParams[accountIndex], 0)
	if !ok {
		s.normalErrorF(rw, -100, "Invalid Account Index")
		return
	}
	// Debug
	logger.Debugf("input : accountID(%v), index(%v)", s.accountID, index)

	// Get address
	priv, err := DefaultWallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	childPk, err := priv.Public().ChildKey(index)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	addr, err := abchainTx.NewAddress(childPk)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, addr.ToString())
}

func (s *account) Update(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received update account request")

	// Parse newAccountID (old ID is still grep from path)
	newAccountID := req.FormValue("newAccountID")
	if newAccountID == "" {
		http.Error(rw, "Missing parameters", http.StatusBadRequest)
	}

	// Debug
	logger.Debugf("input : accountID(%v), newAccountID(%v)", s.accountID, newAccountID)

	// Rename account
	err := DefaultWallet.Rename(s.accountID, newAccountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	// Get address
	priv, err := DefaultWallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	// Persist wallet data
	s.shouldPersist = true

	s.normal(rw, addr.ToString())

}

func (s *account) Delete(rw web.ResponseWriter, req *web.Request) {
	logger.Debug("Received delete account request")

	// Check accountID
	if s.accountID == "" {
		http.Error(rw, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Delete account
	err := DefaultWallet.RemovePrivKey(s.accountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, nil)

	// Persist wallet data
	s.shouldPersist = true
}

func (s *account) ExportKey(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received export private key request")

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Load private key
	priv, err := DefaultWallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.normalError(rw, err)
		return
	}

	s.normal(rw, priv.Str())
}

func (s *account) ImportKey(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received import private key request")

	// Parse privkey
	privkeyS, ok := req.PostForm["privkey"]
	if !ok || len(privkeyS) == 0 {
		http.Error(rw, "Missing parameters", http.StatusBadRequest)
	}

	// Debug
	logger.Debugf("input : accountID(%v), privkey(%v)", s.accountID, privkeyS)

	retAddr := make([]string, 0, len(privkeyS))

	for i, privstr := range privkeyS {
		priv, err := crypto.PrivatekeyFromString(privstr)
		if err != nil {
			s.normalError(rw, err)
			return
		}

		var id string
		if len(privkeyS) == 1 {
			id = s.accountID
		} else {
			id = fmt.Sprintf("%s_%d", s.accountID, i)
		}

		// Import private key
		err = DefaultWallet.ImportPrivateKey(id, priv)
		if err != nil {
			s.normalError(rw, err)
			return
		}

		// Generate Address
		addr, err := abchainTx.NewAddressFromPrivateKey(priv)
		if err != nil {
			s.normalError(rw, err)
			return
		}

		retAddr = append(retAddr, addr.ToString())
	}

	// Persist wallet data
	s.shouldPersist = true

	if len(retAddr) > 1 {
		s.normal(rw, retAddr)
	} else {
		s.normal(rw, retAddr[0])
	}

}
