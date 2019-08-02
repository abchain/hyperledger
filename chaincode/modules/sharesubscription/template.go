package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralInvokingTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()
	ret[Method_NewContract] = &tx.ChaincodeTx{ccname, NewContractHandler(cfg), nil, nil}
	ret[Method_Redeem] = &tx.ChaincodeTx{ccname, RedeemHandler(cfg),
		[]tx.TxPreHandler{NewRedeemContractAddrCred(cfg)}, nil}

	return
}

func GeneralQueryTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	ret[Method_MemberQuery] = &tx.ChaincodeTx{ccname, MemberQueryHandler(cfg), nil, nil}
	ret[Method_Query] = &tx.ChaincodeTx{ccname, QueryHandler(cfg), nil, nil}
	return
}

//extend credential for delegator in each contract:
//a post handler and addr verifier is added for newcontract method
//a verifier is added for query method
//credential is required for memberquery method
func ExtendTemplateForDelegator(cts tx.CollectiveTxs, cfg *StandardContractConfig) tx.CollectiveTxs {

	//newcontract is supposed to be existed, or just panic
	cH := cts[Method_NewContract]
	cH.PostHandlers = append(cH.PostHandlers, NewContractVerifier(cfg))
	cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(nil))

	la := GetDeletagorAddress(cfg)
	cH = cts[Method_Redeem]
	cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(la))

	if cH, ok := cts[Method_Query]; ok {
		cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(la))
	}

	if cH, ok := cts[Method_MemberQuery]; ok {
		cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(la))
	}

	return cts
}
