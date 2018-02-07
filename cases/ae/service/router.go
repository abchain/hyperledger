package service

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/crypto"
	"hyperledger.abchain.org/utils"
	"math/big"
	"net/http"
)

var URIPrefix = "/api/v1/"

type serviceCore struct {
}

type apiCore struct {
	*serviceCore
	debugData   interface{}
	activePrivk *crypto.PrivateKey
}

func (s *serviceCore) NotFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%v Not Found", r.URL.Path)
}

//always try to parse form and get a privatekey from wallet
func (s *apiCore) PrehandlePost(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {
	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		//should allow error or ID is not provided
		privk, err := DefaultWallet.LoadPrivKey(req.PostFormValue(accountID))
		if err == nil {
			index, ok := big.NewInt(0).SetString(req.PostFormValue(accountIndex), 0)
			if ok {
				privk, err = privk.ChildKey(index)
				if err != nil {
					s.normalError(rw, err)
					return
				}
			}

			s.activePrivk = privk
		}
	}

	next(rw, req)
}

func (s *apiCore) normalHeader(rw web.ResponseWriter) {

	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS (default option handler will handle OPTION and set Access-Control-Allow-Method properly)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	// Set response status ok
	rw.WriteHeader(http.StatusOK)

}

func (s *apiCore) normal(rw web.ResponseWriter, v interface{}) {

	s.normalHeader(rw)
	// Create response encoder
	json.NewEncoder(rw).Encode(utils.JRPCSuccess(v))
}

func (s *apiCore) normalError(rw web.ResponseWriter, e error) {

	s.normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCError(e, s.debugData))
}

func (s *apiCore) normalErrorF(rw web.ResponseWriter, code int, message string) {

	s.normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCErrorF(code, message, s.debugData))
}

func buildRouter() *web.Router {

	root := web.New(serviceCore{})

	router := root.Subrouter(apiCore{}, URIPrefix)
	router.Middleware((*apiCore).PrehandlePost)

	// Deploy
	depRouter := router.Subrouter(deploy{}, "deploy")
	depRouter.Post("/", (*deploy).Deploy)

	// Account
	accRouter := router.Subrouter(account{shouldPersist: false}, "account")
	accRouter.Middleware((*account).PersistAccount).
		Middleware((*account).ParseParameters)
	accRouter.Post("/", (*account).Create)
	accRouter.Get("/", (*account).List)
	accRouter.Get(accountIDDir, (*account).Query)
	accRouter.Patch(accountIDDir, (*account).Update)
	accRouter.Delete(accountIDDir, (*account).Delete)
	accRouter.Get(accountIDDir+accountIndexDir, (*account).QueryChild)

	// Account - private
	privRouter := router.Subrouter(account{shouldPersist: false}, "privkey")
	privRouter.Middleware((*account).PersistAccount).
		Middleware((*account).ParseParameters)
	privRouter.Post("/", (*account).ImportKey)
	privRouter.Get(accountIDDir, (*account).ExportKey)

	// // Registrar
	regRouter := router.Subrouter(registrar{}, "registrar")
	regRouter.Middleware((*registrar).InitCaller)
	regRouter.Post("/", (*registrar).Reg)
	regRouter.Get(regPkIDDir, (*registrar).Query)
	// regRouter.Post("/audit", (*RegistrarREST).Audit)

	//token assign
	globalRouter := router.Subrouter(fund{}, "assign")
	globalRouter.Middleware((*fund).InitCaller)
	globalRouter.Post("/", (*fund).Assign)
	globalRouter.Get("/", (*fund).QueryGlobal)

	//Fund
	fundRouter := router.Subrouter(fund{}, "fund")
	fundRouter.Middleware((*fund).InitCaller)
	fundRouter.Post("/", (*fund).Fund)
	fundRouter.Get(fundIDDir, (*fund).QueryTransfer)

	//share
	shareRouter := router.Subrouter(subscription{}, "subscription")
	shareRouter.Middleware((*subscription).InitCaller)
	shareRouter.Post("/", (*subscription).NewContract)
	shareRouter.Post("/redeem"+contrcatAddrDir, (*subscription).Redeem)
	shareRouter.Get(contrcatAddrDir, (*subscription).QueryContract)

	// BlockChain - Address
	addresRouter := router.Subrouter(fund{}, "address")
	addresRouter.Middleware((*fund).InitCaller)
	addresRouter.Get(addressFlagDir, (*fund).QueryAddress)
	//addresRouter.Get("/:accountID/:index", (*fund).Query)

	// BlockChain
	// blockchainRouter := router.Subrouter(BlockChainREST{}, "chain")
	// blockchainRouter.Get("/", (*BlockChainREST).QueryChain)
	// blockchainRouter.Get("/blocks/:height", (*BlockChainREST).QueryBlock)
	// blockchainRouter.Get("/transactions/:transactionID", (*BlockChainREST).QueryTransaction)

	// NotFound
	root.NotFound((*serviceCore).NotFound)

	return root
}
