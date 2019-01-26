package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	"github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	"hyperledger.abchain.org/chaincode/lib/caller"
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

type RPCRouter struct {
	*web.Router
}

func CreateRPCRouter(root *web.Router, path string) RPCRouter {
	return RPCRouter{
		root.Subrouter(FabricRPCBase{}, path),
	}
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
