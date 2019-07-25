package service

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/asset/currency"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/applications/supervise/multisig"
	"hyperledger.abchain.org/applications/supervise/registar"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/applications/util/blockchain"
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
	wallet.CreateAccountRouter(root, "account").Init(defaultWallet).BuildRoutes()

	//a temprouter, remove later
	wallet.CreateAccountRouter(root, "").Post("/address", (*wallet.Account).PublicKeyToAddress)
	//privkey
	wallet.CreateAccountRouter(root, "privkey").Init(defaultWallet).BuildPrivkeyRoutes()

	buildBusiness := func(rpc util.RPCRouter) util.TxRouter {

		root := util.CreateTxRouter(rpc).Init(chaincode.CC_NAME)
		root.BuildRoutes()

		//assign
		currency.CreateFundRouter(root, "assign").Init().BuildGlobalRoutes()
		currency.CreateFundRouter(root, currency.TokenNamePath+"/assign").Init().BuildGlobalRoutes()

		//fund
		currency.CreateFundRouter(root, "fund").Init().BuildFundRoutes()
		currency.CreateFundRouter(root, currency.TokenNamePath+"/fund").Init().BuildFundRoutes()

		//address
		currency.CreateFundRouter(root, "address").Init().BuildAddressRoutes()
		currency.CreateFundRouter(root, currency.TokenNamePath+"/address").Init().BuildAddressRoutes()

		//share
		currency.CreatSubscriptionRouter(root, "subscription").Init().BuildRoutes()

		//fundbatch
		batchroot := util.CreateBatchRouter(root, "adv").Init(chaincode.CC_BATCH)
		currency.CreateFundBatchRouter(batchroot, "").BuildRoutes()

		//registrar
		registar.CreatRegistrarRouter(root, "registrar").Init().BuildRoutes()

		//mauth
		multisig.CreatMultisignRouter(root, "mauth").Init().BuildRoutes()

		return root
	}

	apirouter := util.CreateRPCRouter(root, "").Init(defaultRpcCaller)
	localrouter := util.CreateRPCRouter(root, "data").Init(util.MakeDefaultLocalCaller)

	//business
	wallet.InitTxRouterWithWallet(buildBusiness(apirouter), defaultWallet)
	wallet.InitTxRouterWithWallet(buildBusiness(localrouter), defaultWallet)

	//blockchain
	blockchain.CreateBlockChainRouter(apirouter, "chain").Init(defaultChain).BuildRoutes()

	// NotFound
	root.NotFound(notFound)

	return root
}
