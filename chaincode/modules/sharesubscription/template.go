package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralInvokingTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	cH := &tx.ChaincodeTx{ccname, NewContractHandler(cfg), nil, nil}
	cH.PreHandlers = append(cH.PreHandlers, tx.NewAddrCredVerifier(nil))
	ret[Method_NewContract] = cH

	rcH := &tx.ChaincodeTx{ccname, RedeemHandler(cfg), nil, nil}
	rcH.PreHandlers = append(rcH.PreHandlers, NewRedeemContractAddrCred(cfg))

	ret[Method_Redeem] = rcH

	return
}

func GeneralQueryTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	ret[Method_MemberQuery] = &tx.ChaincodeTx{ccname, MemberQueryHandler(cfg), nil, nil}
	ret[Method_Query] = &tx.ChaincodeTx{ccname, QueryHandler(cfg), nil, nil}
	return
}
