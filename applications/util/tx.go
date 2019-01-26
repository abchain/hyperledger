package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
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

func (s *FabricRPCCore) EncodeEntry(nonce []byte) string {
	return base64.URLEncoding.EncodeToString(nonce)
}

func (s *FabricRPCCore) DecodeEntry(nonce string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(nonce)
}

func (s *FabricRPCCore) SendRawTx(rw web.ResponseWriter, req *web.Request) {

	err, flag, method, _, args := ParseCompactFormTx(req.PostFormValue("tx"))

	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txMaker := txutil.NewTxMaker(args)
	sigs := req.PostForm["sig"]

	for i, sig := range sigs {
		sigpb, err := crypto.DecodeCompactSignature(sig)
		if err != nil {
			s.NormalError(rw, fmt.Errorf("Decode signature %d fail: %s", i, err))
			return
		}
		txMaker.GetCredBuilder().AddSignature(sigpb)
	}

	if s.ActivePrivk != nil {
		sig, err := s.ActivePrivk.Sign(txMaker.GenHash(method))
		if err != nil {
			s.NormalError(rw, fmt.Errorf("Sign tx from privatekey fail: %s", err))
			return
		}
		txMaker.GetCredBuilder().AddSignature(sig)
	}

	args, err = txMaker.GenArguments()
	if err != nil {
		s.NormalError(rw, fmt.Errorf("Gen signed args fail: %s", err))
		return
	}

	var retTx string
	switch flag {
	case "I":
		retTx, err = s.Caller.Invoke(method, args)
	case "D":
		retTx, err = s.Caller.Deploy(method, args)
	case "Q":
		s.NormalErrorF(rw, 500, "Not implied yet")
		return
	default:
		s.NormalError(rw, fmt.Errorf("No such a tx type: %s", flag))
		return
	}

	if err != nil {
		s.NormalError(rw, err)
	} else {
		s.Normal(rw, retTx)
	}

}

func (s *FabricRPCCore) DoSignature(rw web.ResponseWriter, req *web.Request) {

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
	}

	hash := req.PostFormValue("hash")
	var hashbt []byte

	if _, err := fmt.Sscanf(hash, "%x", &hashbt); err != nil {
		s.NormalError(rw, err)
		return
	}

	if sig, err := s.ActivePrivk.Sign(hashbt); err != nil {
		s.NormalError(rw, err)
	} else if sigstr, err := crypto.EncodeCompactSignature(sig); err != nil {
		s.NormalError(rw, err)
	} else {
		s.Normal(rw, sigstr)
	}

}

type TxRouter struct {
	*web.Router
}

func CreateTxRouter(root RPCRouter) TxRouter {
	return TxRouter{
		root.Subrouter(FabricRPCCore{}, ""),
	}
}

func (r TxRouter) BuildRoutes() {
	r.Post("/rawtransaction", (*FabricRPCCore).SendRawTx)
	r.Post("/signature", (*FabricRPCCore).DoSignature)
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
