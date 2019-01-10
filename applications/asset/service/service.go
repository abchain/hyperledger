package service

import (
	"math/big"

	"github.com/gocraft/web"
	log "github.com/op/go-logging"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/core/crypto"
)

var logger = log.MustGetLogger("server/asset")

type RPCCoreWithAccount struct {
	*util.FabricRPCCore
	wallet wallet.Wallet
	// ActivePrivk *crypto.PrivateKey
	ActivePrivk crypto.Signer
}

type RPCAccountRouter struct {
	*web.Router
}

func CreateRPCAccountRouter(root util.RPCRouter, path string) RPCAccountRouter {
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
			if indstr := req.PostFormValue(AccountIndex); indstr != "" {
				index, ok := big.NewInt(0).SetString(indstr, 0)
				if ok {
					privk, err = crypto.GetChildPrivateKey(privk, index)
					if err != nil {
						s.NormalError(rw, err)
						return
					}
				}
			}
			s.ActivePrivk = privk
		}

		s.wallet = wallet
		next(rw, req)
	}

	r.Middleware(Initcall)
}
