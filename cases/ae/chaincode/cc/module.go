package chaincode

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	_ "hyperledger.abchain.org/chaincode/modules/registrar"
	share "hyperledger.abchain.org/chaincode/modules/sharesubscription"
)

type AECC struct {
	tx.CollectiveTxs
	runtimeCfg *runtime.Config
}

const (
	CC_NAME = "AtomicEnergy_v1"
	CC_TAG  = "AE"
)

func NewChaincode(debugMode bool) *AECC {

	ret := &AECC{runtimeCfg: runtime.NewConfig()}

	tokencfg := token.NewConfig(CC_TAG)
	tokencfg.Config = ret.runtimeCfg

	handlers := token.GeneralAdminTemplate(CC_NAME, tokencfg)
	handlers = handlers.MustMerge(
		token.GeneralInvokingTemplate(CC_NAME, tokencfg),
		token.GeneralQueryTemplate(CC_NAME, tokencfg),
	)

	sharecfg := share.NewConfig(CC_TAG)
	sharecfg.Config = ret.runtimeCfg

	handlers = handlers.MustMerge(
		share.GeneralInvokingTemplate(CC_NAME, sharecfg),
		share.GeneralQueryTemplate(CC_NAME, sharecfg),
	)

	if !debugMode {
		//build init batch function ...
		initH := token.GeneralAdminTemplate(CC_NAME, tokencfg)
		handlers["init"] = &tx.ChaincodeTx{CC_NAME, tx.BatchTxHandler(initH), nil, nil}
		delete(handlers.Map(), token.Method_Init) //remove this method from general handling, but reserve assign

		//TODO: add verifier for assign ...
	}

	ret.CollectiveTxs = handlers

	return ret
}
