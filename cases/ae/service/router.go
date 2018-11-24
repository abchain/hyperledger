package service

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	mainsrv "hyperledger.abchain.org/applications/asset/service"
	"hyperledger.abchain.org/applications/blockchain"
	regsrv "hyperledger.abchain.org/applications/supervise/service"
	"hyperledger.abchain.org/applications/util"
	"strings"
)

var URIPrefix = "/api/v1/"

func notFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%v Not Found", r.URL.Path)
}

func optionsHandler(rw web.ResponseWriter, r *web.Request, methods []string) {

	rw.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	rw.Header().Add("Access-Control-Allow-Origin", "*")

}

func buildRouter() *web.Router {

	root := web.NewWithPrefix(util.FabricClientBase{}, URIPrefix)

	root.OptionsHandler(optionsHandler)
	//account
	mainsrv.CreateAccountRouter(root, "account").Init(defaultWallet).BuildRoutes()
	//privkey
	mainsrv.CreateAccountRouter(root, "privkey").Init(defaultWallet).BuildPrivkeyRoutes()

	//blockchain
	blockchain.CreateBlocChainRouter(root, "chain").Init(defaultChainConfig).BuildRoutes()

	rpcrouter := util.CreateRPCRouter(root, "")
	rpcrouter.Init(defaultRpcConfig)

	apirouter := mainsrv.CreateRPCAccountRouter(rpcrouter, "")
	apirouter.Init(defaultWallet)

	//assign
	mainsrv.CreateFundRouter(apirouter, "assign").Init().BuildGlobalRoutes()

	//fund
	mainsrv.CreateFundRouter(apirouter, "fund").Init().BuildFundRoutes()

	//address
	mainsrv.CreateFundRouter(apirouter, "address").Init().BuildAddressRoutes()

	//share
	mainsrv.CreatSubscriptionRouter(apirouter, "subscription").Init().BuildRoutes()

	//registrar
	regsrv.CreatRegistrarRouter(apirouter, "registrar").Init().BuildRoutes()

	// NotFound
	root.NotFound(notFound)

	return root
}
