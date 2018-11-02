package service

import (
	"github.com/gocraft/web"
	log "github.com/op/go-logging"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/crypto"
	"math/big"
)

var logger = log.MustGetLogger("server")

type RPCCoreWithAccount struct {
	*client.FabricRPCCore
	wallet      wallet.Wallet
	ActivePrivk *crypto.PrivateKey
}

type RPCAccountRouter struct {
	*web.Router
}

func CreateRPCAccountRouter(root client.RPCRouter, path string) RPCAccountRouter {
	return RPCAccountRouter{
		root.Subrouter(RPCCoreWithAccount{}, path),
	}
}

func (r RPCAccountRouter) Init(wallet wallet.Wallet) {

	Initcall := func(s *RPCCoreWithAccount, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		//should allow error or ID is not provided
		privk, err := wallet.LoadPrivKey(req.PostFormValue(AccountID))
		if err == nil {
			index, ok := big.NewInt(0).SetString(req.PostFormValue(AccountIndex), 0)
			if ok {
				privk, err = privk.ChildKey(index)
				if err != nil {
					s.NormalError(rw, err)
					return
				}
			}

			s.ActivePrivk = privk
		}

		s.wallet = wallet
		next(rw, req)
	}

	r.Middleware(Initcall)
}
