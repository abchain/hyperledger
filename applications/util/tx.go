package util

import (
	//"encoding/base64"
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

	nonceStr    string
	rwForOutput web.ResponseWriter
}

func (s *FabricRPCCore) PrehandlePost(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		s.nonceStr = req.PostFormValue("nonce")
		if s.nonceStr != "" {
			s.TxGenerator.BeginTx([]byte(s.nonceStr))
		} else {
			s.TxGenerator.BeginTx(nil)
		}

		s.rwForOutput = rw
	}

	next(rw, req)
}

type defaultEntry struct {
	Txid  string      `json:"txID"`
	Nonce string      `json:"Nonce"`
	Entry interface{} `json:"Data,omitempty"`
}

func (s *FabricRPCCore) DefaultOutput(e interface{}) {

	if s.rwForOutput == nil {
		//sanity check
		panic("Should not called in non-POST (invoking) method")
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(s.rwForOutput, err)
		return
	}

	var ncStr string
	if s.nonceStr != "" {
		ncStr = fmt.Sprintf("%X (%s)", s.TxGenerator.GetNonce(), s.nonceStr)
	} else {
		ncStr = fmt.Sprintf("%X", s.TxGenerator.GetNonce())
	}

	s.Normal(s.rwForOutput, &defaultEntry{txid, ncStr, e})
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
