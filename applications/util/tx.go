package util

import (
	"fmt"
	"github.com/gocraft/web"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/core/crypto"
	"net/http"
)

type FabricRPCCore struct {
	*FabricRPCBase
	*txgen.TxGenerator
	ActivePrivk crypto.Signer
}

func (s *FabricRPCCore) PrehandlePost(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {
	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		nonce := req.PostFormValue("nonce")
		if nonce != "" {
			s.TxGenerator.BeginTx([]byte(nonce))
		} else {
			s.TxGenerator.BeginTx(nil)
		}
	}

	next(rw, req)
}

type TxRouter struct {
	*web.Router
}

func CreateTxRouter(root RPCRouter) TxRouter {
	return TxRouter{
		root.Subrouter(FabricRPCCore{}, ""),
	}
}

func (r TxRouter) Init(ccname string) TxRouter {

	initCall := func(s *FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		s.TxGenerator = txgen.SimpleTxGen(ccname)
		//detect local caller
		if lc, ok := s.Caller.(LocalCaller); ok {
			s.RespWrapping = func(msg interface{}) interface{} {

				var encoded struct {
					RawTx   string      `json:"raw"`
					Txhash  string      `json:"hash"`
					Promise interface{} `json:"promise,omitempty"`
				}

				encoded.RawTx = lc.Output()
				encoded.Promise = msg
				encoded.Txhash = fmt.Sprintf("%X", s.GetBuilder().GetHash())

				return &encoded
			}

		}

		s.TxGenerator.Dispatcher = s.Caller
		next(rw, req)
	}

	r.Middleware(initCall).
		Middleware((*FabricRPCCore).PrehandlePost)
	return r
}
