package service

import (
	"github.com/gocraft/web"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	reg "hyperledger.abchain.org/chaincode/registrar"
	"hyperledger.abchain.org/client"
	"math/big"
	"net/http"
)

type deploy struct {
	*client.FabricRPCCore
}

func (s *deploy) Deploy(rw web.ResponseWriter, req *web.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var args []string

	//token deployment
	total, ok := big.NewInt(0).SetString(req.PostFormValue("total"), 0)
	if !ok {
		http.Error(rw, "Invalid total parameter", http.StatusBadRequest)
		return
	}

	args, err = token.CCDeploy(total, args)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	manager := req.FormValue("Admin")
	if manager == "" {
		manager = "Admin"
	}
	regmanager := req.FormValue("RegManager")
	if regmanager == "" {
		regmanager = manager
	}
	args, err = reg.CCDeploy(manager, regmanager, args)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !offlineMode {
		caller, err := defaultRpcConfig.NewCall()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		caller.Function = "INIT"
		_, err = caller.Deploy(args)
	} else {
		_, err = ccCaller.MockInit("regtest_deploy_test", "INIT", args)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Normal(rw, "OK")
}
