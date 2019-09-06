package ticket

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralInvokingTemplate(ccname string, cfg TicketConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	cH := &tx.ChaincodeTx{ccname, AddHandler(cfg), nil, nil}
	ret[Method_Add] = cH

	rcH := &tx.ChaincodeTx{ccname, ApplyHandler(cfg), nil, nil}
	ret[Method_Apply] = rcH

	return
}

func GeneralQueryTemplate(ccname string, cfg TicketConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()
	ret[Method_Query] = &tx.ChaincodeTx{ccname, QueryHandler(cfg), nil, nil}
	return
}
