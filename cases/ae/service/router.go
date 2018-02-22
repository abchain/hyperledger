package service

import (
	"fmt"
	"github.com/gocraft/web"
	mainsrv "hyperledger.abchain.org/asset/service"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	regsrv "hyperledger.abchain.org/supervise/service"
	"net/http"
)

type rpcCfg struct{}

func (*rpcCfg) GetCCName() string {
	return chaincode.CC_NAME
}

func (*rpcCfg) GetCaller() rpc.Caller {
	if offlineMode {
		return ccCaller
	} else {
		c, err := defaultRpcConfig.NewCall()
		if err != nil {
			return nil
		}

		return c
	}
}

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

	rpcrouter := client.CreateRPCRouter(root, URIPrefix)
	rpcrouter.Init(&rpcCfg{})

	// Deploy
	rpcrouter.Subrouter(deploy{}, "deploy").Post("/", (*deploy).Deploy)

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
