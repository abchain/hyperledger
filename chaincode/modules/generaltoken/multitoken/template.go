package multitoken

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

//admintemplate has no verifier
func GeneralAdminTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {
	ret = tx.NewCollectiveTxs()

	ret[Method_Init] = &tx.ChaincodeTx{ccname, InitHandler(cfg), nil, nil}
	ret[Method_Assign] = &tx.ChaincodeTx{ccname, AssignHandler(cfg), nil, nil}

	return
}

func GeneralInvokingTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	h := TransferHandler(cfg)
	txh := &tx.ChaincodeTx{ccname, h, nil, nil}
	txh.PreHandlers = append(txh.PreHandlers, tx.NewAddrCredVerifier(nil))
	ret[Method_Transfer] = txh
	return
}

func ExtendInvokingTemplate(cts tx.CollectiveTxs, verifier tx.AddrVerifier) tx.CollectiveTxs {

	if h, ok := cts[Method_Transfer]; ok {
		tx.AttachAddrVerifier(h.PreHandlers, verifier)
	}

	return cts
}

func GeneralQueryTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	ret[Method_QueryToken] = &tx.ChaincodeTx{ccname, TokenQueryHandler(cfg), nil, nil}
	ret[Method_QueryGlobal] = &tx.ChaincodeTx{ccname, GlobalQueryHandler(cfg), nil, nil}
	return
}
