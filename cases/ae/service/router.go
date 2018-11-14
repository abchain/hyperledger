package service

import (
	"fmt"
	"github.com/gocraft/web"
	mainsrv "hyperledger.abchain.org/applications/asset/service"
	regsrv "hyperledger.abchain.org/applications/supervise/service"
	"hyperledger.abchain.org/client"
	"net/http"
)

var URIPrefix = "/api/v1/"

func notFound(w web.ResponseWriter, r *web.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%v Not Found", r.URL.Path)
}
func buildRouter() *web.Router {

	root := web.NewWithPrefix(client.FabricClientBase{}, URIPrefix)

	//account
	mainsrv.CreateAccountRouter(root, "account").Init(defaultWallet).BuildRoutes()
	//privkey
	mainsrv.CreateAccountRouter(root, "privkey").Init(defaultWallet).BuildPrivkeyRoutes()

	//blockchain
	client.CreateFabricProxyRouter(root, "chain").Init(defaultFabricEP, nil).BuildRoutes()

	rpcrouter := client.CreateRPCRouter(root, "")
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
