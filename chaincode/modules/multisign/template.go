package multisign

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralInvokingTemplate(ccname string, cfg MultiSignConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	verifier := MultiSignAddrPreHandler(cfg)

	cH := &tx.ChaincodeTx{ccname, ContractHandler(cfg), nil, nil}
	ret[Method_Contract] = cH

	rcH := &tx.ChaincodeTx{ccname, UpdateHandler(cfg), nil, nil}
	rcH.PreHandlers = append(rcH.PreHandlers, verifier)
	ret[Method_Update] = rcH

	return
}

func GeneralQueryTemplate(ccname string, cfg MultiSignConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()
	ret[Method_Query] = &tx.ChaincodeTx{ccname, QueryHandler(cfg), nil, nil}
	return
}
