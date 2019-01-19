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

	rpcrouter := util.CreateRPCRouter(root).Init(defaultRpcConfig)
	rpcrouter.BuildRoutes()

	//account
	mainsrv.CreateAccountRouter(root, "account").Init(defaultWallet).BuildRoutes()
	//privkey
	mainsrv.CreateAccountRouter(root, "privkey").Init(defaultWallet).BuildPrivkeyRoutes()

	//blockchain
	blockchain.CreateBlocChainRouter(rpcrouter, "chain").Init(defaultChainConfig).BuildRoutes()

	apirouter := util.CreateTxRouter(rpcrouter, "").Init(defaultRpcConfig.GetCCName())
	localrouter := util.CreateTxRouter(rpcrouter, "data").InitLocalCall(defaultRpcConfig.GetCCName())

	buildBusiness := func(root util.TxRouter) {

		//assign
		mainsrv.CreateFundRouter(root, "assign").Init().BuildGlobalRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/assign").Init().BuildGlobalRoutes()

		//fund
		mainsrv.CreateFundRouter(root, "fund").Init().BuildFundRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/fund").Init().BuildFundRoutes()

		//address
		mainsrv.CreateFundRouter(root, "address").Init().BuildAddressRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/address").Init().BuildAddressRoutes()

		//share
		mainsrv.CreatSubscriptionRouter(root, "subscription").Init().BuildRoutes()

		//registrar
		regsrv.CreatRegistrarRouter(root, "registrar").Init().BuildRoutes()

	}

	mainsrv.InitTxRouterWithWallet(apirouter, defaultWallet)
	buildBusiness(apirouter)
	buildBusiness(localrouter)

	// NotFound
	root.NotFound(notFound)

	return root
}
