package util

import (
	"encoding/base64"
	"fmt"
	"github.com/gocraft/web"
	"github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
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

type localCaller struct {
	txutil.Builder
	encodedTx struct {
		RawTx   string      `json:"raw"`
		Txhash  string      `json:"hash"`
		Promise interface{} `json:"promise,omitempty"`
	}
}

func encodeArguments(args [][]byte) (r []string) {
	for _, arg := range args {
		r = append(r, base64.StdEncoding.EncodeToString(arg))
	}

	return
}

func (c *localCaller) initResp(txType string, method string, arg [][]byte) {
	c.encodedTx.RawTx = strings.Join(append([]string{txType, method}, encodeArguments(arg)...), ":")
	c.encodedTx.Txhash = fmt.Sprintf("%X", c.GetHash())
}

func (c *localCaller) Deploy(method string, arg [][]byte) (string, error) {

	return "pending", nil
}

func (c *localCaller) Invoke(method string, arg [][]byte) (string, error) {
	c.encodedTx.RawTx = strings.Join(append([]string{"I", method}, encodeArguments(arg)...), ":")
	return "pending", nil
}

func (c *localCaller) Query(method string, arg [][]byte) ([]byte, error) {
	c.encodedTx.RawTx = strings.Join(append([]string{"Q", method}, encodeArguments(arg)...), ":")
	return emptyProto, nil
}

func (c *localCaller) Wrapping(msg interface{}) interface{} {
	c.encodedTx.Promise = msg
	return &c.encodedTx
}

type LocalTxDataRouter struct {
	*web.Router
}

func CreateLocalTxRouter(root *web.Router, path string) LocalTxDataRouter {
	return LocalTxDataRouter{
		root.Subrouter(FabricRPCCore{}, path),
	}
}

func (r LocalTxDataRouter) Init(ccname string) {

	initCall := func(s *FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {
		s.TxGenerator = txgen.SimpleTxGen(ccname)
		lc := &localCaller{}
		lc.Builder = s.GetBuilder()
		s.TxGenerator.Dispatcher = lc

		next(rw, req)
	}

	r.Middleware(initCall).
		Middleware((*FabricRPCCore).PrehandlePost)
}
