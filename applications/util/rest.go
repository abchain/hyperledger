package util

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/core/utils"
	"net/http"
)

//a null-base to provide more elastic
type FabricClientBase struct {
	debugData interface{}
}

//so we can support both config in client and local adapter
type FabricRPCCfg interface {
	GetCaller() (rpc.Caller, error)
	GetCCName() string
	Quit()
}

type FabricRPCCore struct {
	*FabricClientBase
	*txgen.TxGenerator
	Cfg FabricRPCCfg
}

type RPCRouter struct {
	*web.Router
}

func CreateRPCRouter(root *web.Router, path string) RPCRouter {
	return RPCRouter{
		root.Subrouter(FabricRPCCore{}, path),
	}
}

func (r RPCRouter) Init(cfg FabricRPCCfg) {

	initCall := func(s *FabricRPCCore, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		var err error
		s.TxGenerator = txgen.SimpleTxGen(cfg.GetCCName())
		s.TxGenerator.Dispatcher, err = cfg.GetCaller()

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		} else {
			next(rw, req)
		}
	}

	r.Middleware(initCall).
		Middleware((*FabricRPCCore).PrehandlePost)
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

func (s *FabricClientBase) normalHeader(rw web.ResponseWriter) {

	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS (default option handler will handle OPTION and set Access-Control-Allow-Method properly)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type")

	// Set response status ok
	rw.WriteHeader(http.StatusOK)

}

func (s *FabricClientBase) Normal(rw web.ResponseWriter, v interface{}) {

	s.normalHeader(rw)
	// Create response encoder
	json.NewEncoder(rw).Encode(utils.JRPCSuccess(v))
}

func (s *FabricClientBase) NormalError(rw web.ResponseWriter, e error) {

	s.normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCError(e, s.debugData))
}

func (s *FabricClientBase) NormalErrorF(rw web.ResponseWriter, code int, message string) {

	s.normalHeader(rw)
	json.NewEncoder(rw).Encode(utils.JRPCErrorF(code, message, s.debugData))
}

func (s *FabricClientBase) EncodeEntry(nonce []byte) string {
	return base64.URLEncoding.EncodeToString(nonce)
}

func (s *FabricClientBase) DecodeEntry(nonce string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(nonce)
}
