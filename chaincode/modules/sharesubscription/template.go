package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func SimpleInvokingTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()
	ret[Method_NewContract] = &tx.ChaincodeTx{ccname, NewContractHandler(cfg), nil, nil}
	ret[Method_Redeem] = &tx.ChaincodeTx{ccname, RedeemHandler(cfg), nil, nil}

	return
}

func GeneralInvokingTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = SimpleInvokingTemplate(ccname, cfg)

	rcH := ret[Method_Redeem]
	rcH.PreHandlers = append(rcH.PreHandlers,
		NewRedeemContractAddrCred(cfg),
		tx.NewAddrCredVerifier(nil))

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

	handler := NewContractVerifier(cfg)

	//newcontract is supposed to be existed, or just panic
	cH := cts[Method_NewContract]
	cH.PostHandlers = append(cH.PostHandlers, handler)
	cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(nil))

	if cH, ok := cts[Method_Query]; ok {
		cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(handler.La()))
	}

	if cH, ok := cts[Method_MemberQuery]; ok {
		cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(nil))
	}

	return cts
}
