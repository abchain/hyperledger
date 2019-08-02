package chaincode

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	mtoken "hyperledger.abchain.org/chaincode/modules/generaltoken/multitoken"
	tokenNc "hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	multisign "hyperledger.abchain.org/chaincode/modules/multisign"
	_ "hyperledger.abchain.org/chaincode/modules/registrar"
	share "hyperledger.abchain.org/chaincode/modules/sharesubscription"
	_ "hyperledger.abchain.org/core/crypto/ecdsa" //important
)

type AECC struct {
	tx.CollectiveTxs
	runtimeCfg *runtime.Config
}

const (
	CC_NAME  = "AtomicEnergy_v1"
	CC_TAG   = "AE"
	CC_BATCH = "batch"
)

func NewChaincode(debugMode bool) *AECC {

	ret := &AECC{runtimeCfg: runtime.NewConfig()}

	tokencfg := token.NewConfig(CC_TAG)
	tokencfg.Config = ret.runtimeCfg
	tokenNccfg := tokenNc.NewConfig(CC_TAG)
	tokenNccfg.Config = ret.runtimeCfg

	handlers := token.GeneralAdminTemplate(CC_NAME, tokencfg)

	mauthcfg := multisign.NewConfig(CC_TAG)
	handlers = handlers.MustMerge(multisign.GeneralInvokingTemplate(CC_NAME, mauthcfg),
		multisign.GeneralQueryTemplate(CC_NAME, mauthcfg))

	handlers = handlers.MustMerge(
		token.ExtendInvokingTemplate(token.GeneralInvokingTemplate(CC_NAME, tokencfg),
			multisign.MultiSignAddrPreHandler(mauthcfg)),
		token.GeneralQueryTemplate(CC_NAME, tokencfg),
		tokenNc.GeneralTemplate(CC_NAME, tokenNccfg),
	)

	mtokencfg := mtoken.ConfigFromToken(tokencfg)
	handlers = handlers.MustMerge(
		mtoken.GeneralAdminTemplate(CC_NAME, mtokencfg),
		mtoken.ExtendInvokingTemplate(mtoken.GeneralInvokingTemplate(CC_NAME, mtokencfg),
			multisign.MultiSignAddrPreHandler(mauthcfg)),
		mtoken.GeneralQueryTemplate(CC_NAME, mtokencfg),
	)

	sharecfg := share.NewConfig(CC_TAG)
	sharecfg.Config = ret.runtimeCfg

	handlers = handlers.MustMerge(
		share.GeneralInvokingTemplate(CC_NAME, sharecfg),
		share.GeneralQueryTemplate(CC_NAME, sharecfg),
	)
	handlers = share.ExtendTemplateForDelegator(handlers, sharecfg)

	if !debugMode {
		// //build init batch function ...
		// initH := token.GeneralAdminTemplate(CC_NAME, tokencfg)
		// handlers["init"] = &tx.ChaincodeTx{CC_NAME, tx.BatchTxHandler(initH), nil, nil}
		// delete(handlers.Map(), token.Method_Init) //remove this method from general handling, but reserve assign

		//TODO: add verifier for assign ...
	}

	//allow all methods use batch ...
	handlers[CC_BATCH] = &tx.ChaincodeTx{CC_NAME, tx.BatchTxHandler(handlers.Map()), nil, nil}

	ret.CollectiveTxs = handlers

	return ret
}
