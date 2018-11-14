package service

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/tx"
	"math/big"
	"net/http"
)

const (
	AccountID    = "accountID"
	AccountIndex = "index"
)

type Account struct {
	*util.FabricClientBase
	accountID     string
	wallet        wallet.Wallet
	shouldPersist bool
}

type AccountRouter struct {
	*web.Router
}

func CreateAccountRouter(root *web.Router, path string) AccountRouter {
	return AccountRouter{
		root.Subrouter(Account{}, path),
	}
}

func (r AccountRouter) Init(wallet wallet.Wallet) AccountRouter {

	Initcall := func(s *Account, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		s.wallet = wallet
		s.shouldPersist = false

		next(rw, req)

		if s.shouldPersist {
			wallet.Persist()
		}
	}

	r.Middleware(Initcall).
		Middleware((*Account).ParseParameters)

	return r

}

func (r AccountRouter) BuildRoutes() {
	r.Post("/", (*Account).Create)
	r.Get("/", (*Account).List)
	r.Get("/:"+AccountID, (*Account).Query)
	r.Patch("/:"+AccountID, (*Account).Update)
	r.Delete("/:"+AccountID, (*Account).Delete)
	r.Get("/:"+AccountID+"/:"+AccountIndex, (*Account).QueryChild)
}

func (r AccountRouter) BuildPrivkeyRoutes() {
	r.Post("/", (*Account).ImportKey)
	r.Get("/:"+AccountID, (*Account).ExportKey)
}

func (s *Account) ParseParameters(rw web.ResponseWriter,
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
		s.accountID = req.PostFormValue(AccountID)
	} else {
		s.accountID = req.PathParams[AccountID]
	}

	next(rw, req)
}

func (s *Account) Create(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create account request")

	// Check accountID
	if s.accountID == "" {
		s.NormalError(rw, errors.New("Must provide accountID"))
		return
	}

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Create private key
	priv, err := s.wallet.NewPrivKey(s.accountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	// Debug
	logger.Debugf("output: address(%v)", addr)

	s.Normal(rw, addr.ToString())

	// Persist wallet data
	s.shouldPersist = true
}

func (s *Account) List(rw web.ResponseWriter, req *web.Request) {
	logger.Debug("Received list accounts request")

	// get account list
	privkeys, err := s.wallet.ListAll()
	if err != nil {
		s.NormalError(rw, err)
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

	s.Normal(rw, ret)
}

func (s *Account) Query(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received query account request")

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Get address
	priv, err := s.wallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.NormalErrorF(rw, 404, "account Not Found")
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, addr.ToString())
}

func (s *Account) QueryChild(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received query child account request")

	// Parse index
	index, ok := big.NewInt(0).SetString(req.PathParams[AccountIndex], 0)
	if !ok {
		s.NormalErrorF(rw, -100, "Invalid Account Index")
		return
	}
	// Debug
	logger.Debugf("input : accountID(%v), index(%v)", s.accountID, index)

	// Get address
	priv, err := s.wallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	childPk, err := priv.Public().ChildKey(index)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	addr, err := abchainTx.NewAddress(childPk)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, addr.ToString())
}

func (s *Account) Update(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received update account request")

	// Parse newAccountID (old ID is still grep from path)
	newAccountID := req.FormValue("newAccountID")
	if newAccountID == "" {
		http.Error(rw, "Missing parameters", http.StatusBadRequest)
	}

	// Debug
	logger.Debugf("input : accountID(%v), newAccountID(%v)", s.accountID, newAccountID)

	// Rename account
	err := s.wallet.Rename(s.accountID, newAccountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	// Get address
	priv, err := s.wallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	// Generate Address
	addr, err := abchainTx.NewAddressFromPrivateKey(priv)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	// Persist wallet data
	s.shouldPersist = true

	s.Normal(rw, addr.ToString())

}

func (s *Account) Delete(rw web.ResponseWriter, req *web.Request) {
	logger.Debug("Received delete account request")

	// Check accountID
	if s.accountID == "" {
		http.Error(rw, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Delete account
	err := s.wallet.RemovePrivKey(s.accountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, nil)

	// Persist wallet data
	s.shouldPersist = true
}

func (s *Account) ExportKey(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received export private key request")

	// Debug
	logger.Debugf("input : accountID(%v)", s.accountID)

	// Load private key
	priv, err := s.wallet.LoadPrivKey(s.accountID)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, priv.Str())
}

func (s *Account) ImportKey(rw web.ResponseWriter, req *web.Request) {

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
			s.NormalError(rw, err)
			return
		}

		var id string
		if len(privkeyS) == 1 {
			id = s.accountID
		} else {
			id = fmt.Sprintf("%s_%d", s.accountID, i)
		}

		// Import private key
		err = s.wallet.ImportPrivateKey(id, priv)
		if err != nil {
			s.NormalError(rw, err)
			return
		}

		// Generate Address
		addr, err := abchainTx.NewAddressFromPrivateKey(priv)
		if err != nil {
			s.NormalError(rw, err)
			return
		}

		retAddr = append(retAddr, addr.ToString())
	}

	// Persist wallet data
	s.shouldPersist = true

	if len(retAddr) > 1 {
		s.Normal(rw, retAddr)
	} else {
		s.Normal(rw, retAddr[0])
	}

}
