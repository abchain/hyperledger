package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	"github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
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

func ParseCompactFormTx(ctxstr string) (e error, flag string, method string, cc string, args [][]byte) {
	rawTx := strings.Split(ctxstr, ":")
	if len(rawTx) < 5 {
		e = fmt.Errorf("Invalid tx (only %d part)", len(rawTx))
		return
	}

	args, e = decodeArguments(rawTx[3:])
	if e != nil {
		return
	}

	flag = rawTx[0]
	cc = rawTx[1]
	method = rawTx[2]

	return
}

type LocalCaller interface {
	Output() string
}

type localCaller string

//specified the fabric-tier chaincode name, or "" to use the default one
func MakeLocalCaller(cc string) (rpc.Caller, error) {
	lc := localCaller("")
	return &lc, nil
}

func MakeDefaultLocalCaller() (rpc.Caller, error) {
	return MakeLocalCaller("")
}

func (c *localCaller) initResp(txType string, method string, arg [][]byte) {
	//TODO: we left a space for cc name
	*c = localCaller(strings.Join(append([]string{txType, "", method}, encodeArguments(arg)...), ":"))
}

func (c *localCaller) Output() string { return string(*c) }

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

type FabricRPCBase struct {
	*FabricClientBase
	Caller rpc.Caller
}

func (s *FabricRPCBase) GetAddress(rw web.ResponseWriter, req *web.Request) {
	pkstr := req.PostFormValue("pubkeybuffer")

	pk, err := crypto.DecodeCompactPublicKey(pkstr)
	if err != nil {
		s.NormalError(rw, fmt.Errorf("decode public key fail: %s", err))
		return
	}

	addr, err := txutil.NewAddress(pk)
	if err != nil {
		s.NormalError(rw, fmt.Errorf("create addr fail: %s", err))
		return
	}

	s.Normal(rw, addr.ToString())
}

func (s *FabricRPCBase) SendRawTx(rw web.ResponseWriter, req *web.Request) {

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

type RPCRouter struct {
	*web.Router
}

func CreateRPCRouter(root *web.Router, path string) RPCRouter {
	return RPCRouter{
		root.Subrouter(FabricRPCBase{}, path),
	}
}

func (r RPCRouter) BuildRoutes() {
	r.Post("/rawtransaction", (*FabricRPCBase).SendRawTx)
	r.Post("/address", (*FabricRPCBase).GetAddress)
}

func (r RPCRouter) Init(cf func() (rpc.Caller, error)) RPCRouter {

	initCall := func(s *FabricRPCBase, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		var err error
		s.Caller, err = cf()
		if err != nil {
			s.NormalError(rw, err)
			return
		}
		next(rw, req)
	}

	r.Middleware(initCall)
	return r
}
