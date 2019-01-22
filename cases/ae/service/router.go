package service

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	mainsrv "hyperledger.abchain.org/applications/asset/service"
	"hyperledger.abchain.org/applications/blockchain"
	regsrv "hyperledger.abchain.org/applications/supervise/service"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
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

	buildBusiness := func(rpc util.RPCRouter) util.TxRouter {

		root := util.CreateTxRouter(rpc).Init(chaincode.CC_NAME)

		//assign
		mainsrv.CreateFundRouter(root, "assign").Init().BuildGlobalRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/assign").Init().BuildGlobalRoutes()

		//fund
		mainsrv.CreateFundRouter(root, "fund").Init().BuildFundRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/fund").Init().BuildFundRoutes()

		//fundbatch
		batchroot := util.CreateBatchRouter(root, "adv").Init(chaincode.CC_BATCH)
		mainsrv.CreateFundBatchRouter(batchroot, "").BuildRoutes()

		//address
		mainsrv.CreateFundRouter(root, "address").Init().BuildAddressRoutes()
		mainsrv.CreateFundRouter(root, mainsrv.TokenNamePath+"/address").Init().BuildAddressRoutes()

		//share
		mainsrv.CreatSubscriptionRouter(root, "subscription").Init().BuildRoutes()

		//registrar
		regsrv.CreatRegistrarRouter(root, "registrar").Init().BuildRoutes()

		return root
	}

	apirouter := util.CreateRPCRouter(root, "").Init(defaultRpcCaller)
	localrouter := util.CreateRPCRouter(root, "data").Init(util.MakeDefaultLocalCaller)

	//business
	apirouter.BuildRoutes()

	mainsrv.InitTxRouterWithWallet(buildBusiness(apirouter), defaultWallet)
	buildBusiness(localrouter)

	//blockchain
	blockchain.CreateBlocChainRouter(apirouter, "chain").Init(defaultChain).BuildRoutes()

	// NotFound
	root.NotFound(notFound)

	return root
}
