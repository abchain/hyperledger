package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	"github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
	"net/http"
	"strings"
)

var emptyProto []byte

func init() {
	empty := new(empty.Empty)
	var err error
	emptyProto, err = proto.Marshal(empty)
	if err != nil {
		panic(err)
	}
}

func encodeArguments(args [][]byte) (r []string) {
	for _, arg := range args {
		r = append(r, base64.StdEncoding.EncodeToString(arg))
	}

	return
}

func decodeArguments(args []string) (r [][]byte, e error) {
	for _, arg := range args {

		var br []byte
		br, e = base64.StdEncoding.DecodeString(arg)
		if e != nil {
			return
		}

		r = append(r, br)
	}

	return
}

func (s *FabricRPCBase) SendRawTx(rw web.ResponseWriter, req *web.Request) {

	rawTx := strings.Split(req.PostFormValue("tx"), ":")
	if len(rawTx) < 5 {
		s.NormalErrorF(rw, -100, "Invalid tx")
		return
	}

	caller, err := s.GetCaller(rawTx[1])
	if err != nil {
		s.NormalError(rw, fmt.Errorf("Get caller fail: %s", err))
		return
	}

	args, err := decodeArguments(rawTx[3:])
	if err != nil {
		s.NormalError(rw, fmt.Errorf("Decode arguments fail: %s", err))
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

	args, err = txMaker.GenArguments()
	if err != nil {
		s.NormalError(rw, fmt.Errorf("Gen signed args fail: %s", err))
		return
	}

	var retTx string
	switch rawTx[0] {
	case "I":
		retTx, err = caller.Invoke(rawTx[2], args)
	case "D":
		retTx, err = caller.Deploy(rawTx[2], args)
	case "Q":
		s.NormalErrorF(rw, 500, "Not implied yet")
		return
	default:
		s.NormalError(rw, fmt.Errorf("No such a tx type: %s", rawTx[0]))
		return
	}

	if err != nil {
		s.NormalError(rw, err)
	} else {
		s.Normal(rw, retTx)
	}

}

type localCaller string

func (c *localCaller) initResp(txType string, method string, arg [][]byte) {
	//TODO: we left a space for cc name
	*c = localCaller(strings.Join(append([]string{txType, "", method}, encodeArguments(arg)...), ":"))
}

func (c *localCaller) Deploy(method string, arg [][]byte) (string, error) {

	c.initResp("D", method, arg)
	return "pending", nil
}

func (c *localCaller) Invoke(method string, arg [][]byte) (string, error) {
	c.initResp("I", method, arg)
	return "pending", nil
}

func (c *localCaller) Query(method string, arg [][]byte) ([]byte, error) {
	c.initResp("Q", method, arg)
	return emptyProto, nil
}

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

func CreateTxRouter(root RPCRouter, path string) TxRouter {
	return TxRouter{
		root.Subrouter(FabricRPCCore{}, path),
	}
}

func (r TxRouter) Init(ccname string) TxRouter {

	initCall := func(s *FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		var err error
		s.TxGenerator = txgen.SimpleTxGen(ccname)
		s.TxGenerator.Dispatcher, err = s.GetCaller(ccname)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			next(rw, req)
		}
	}

	r.Middleware(initCall).
		Middleware((*FabricRPCCore).PrehandlePost)
	return r
}

func (r TxRouter) InitLocalCall(ccname string) TxRouter {

	initCall := func(s *FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {
		s.TxGenerator = txgen.SimpleTxGen(ccname)
		lc := localCaller("")
		s.TxGenerator.Dispatcher = &lc

		s.RespWrapping = func(msg interface{}) interface{} {

			var encoded struct {
				RawTx   string      `json:"raw"`
				Txhash  string      `json:"hash"`
				Promise interface{} `json:"promise,omitempty"`
			}

			encoded.RawTx = string(lc)
			encoded.Promise = msg
			encoded.Txhash = fmt.Sprintf("%X", s.GetBuilder().GetHash())

			return &encoded
		}

		next(rw, req)
	}

	r.Middleware(initCall).
		Middleware((*FabricRPCCore).PrehandlePost)
	return r
}
