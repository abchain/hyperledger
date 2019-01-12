package service

import (
	"math/big"

	"github.com/gocraft/web"
	log "github.com/op/go-logging"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/applications/util"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
)

var logger = log.MustGetLogger("server/asset")

type RPCAccountRouter struct {
	*web.Router
}

func CreateRPCAccountRouter(root util.RPCRouter, path string) RPCAccountRouter {
	return RPCAccountRouter{
		root.Subrouter(util.FabricRPCCore{}, path),
	}
}

func (r RPCAccountRouter) Init(wallet wallet.Wallet) {

	Initcall := func(s *util.FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		//should allow error or ID is not provided
		privk, err := wallet.LoadPrivKey(req.FormValue(AccountID))
		if err == nil {
			if indstr := req.FormValue(AccountIndex); indstr != "" {
				index, ok := big.NewInt(0).SetString(indstr, 0)
				if ok {
					privk, err = crypto.GetChildPrivateKey(privk, index)
					if err != nil {
						s.NormalError(rw, err)
						return
					}
				}
			}
			s.Credgenerator = txgen.NewSingleKeyCred(privk)
			s.ActivePrivk = privk
		}

		next(rw, req)
	}

	r.Middleware(Initcall)
}
