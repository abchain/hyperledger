package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/core/utils"
)

//a null-base to provide more elastic
type FabricClientBase struct {
	debugData    interface{}
	RespWrapping func(interface{}) interface{}
}

type FabricRPCBase struct {
	*FabricClientBase
	GetCaller func(string) (rpc.Caller, error)
}

type RPCRouter struct {
	*web.Router
}

func CreateRPCRouter(root *web.Router) RPCRouter {
	return RPCRouter{
		root.Subrouter(FabricRPCBase{}, ""),
	}
}

func (r RPCRouter) BuildRoutes() {
	r.Post("/sendrawtransaction", (*FabricRPCBase).SendRawTx)
}

//so we can support both config in client and local adapter
type FabricRPCCfg interface {
	GetCaller() (rpc.Caller, error)
	GetCCName() string
	Quit()
}

func (r RPCRouter) Init(cfg FabricRPCCfg) RPCRouter {

	ccn := cfg.GetCCName()
	f := func(n string) (rpc.Caller, error) {
		if n != "" && n != ccn {
			return nil, fmt.Errorf("CC is not match (expect %s but has %s)", ccn, n)
		}

		return cfg.GetCaller()
	}

	initCall := func(s *FabricRPCBase, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		s.GetCaller = f
		next(rw, req)
	}

	r.Middleware(initCall)
	return r
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
	if s.RespWrapping != nil {
		v = s.RespWrapping(v)
	}

	logger.Debugf("Normal finish, output %v", v)

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
