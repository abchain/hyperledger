package generaltoken

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
)

//admintemplate has no verifier
func GeneralAdminTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {
	ret = tx.NewCollectiveTxs()

	ret[Method_Init] = &tx.ChaincodeTx{ccname, InitHandler(cfg), nil, nil}
	ret[Method_Assign] = &tx.ChaincodeTx{ccname, AssignHandler(cfg), nil, nil}

	return
}

func InnerInvokingTemplate(ccname string, cfg *StandardTokenConfig) (ret tx.CollectiveTxs) {
	ret = tx.NewCollectiveTxs()
	ib := &tx.InnerAddrBase{cfg.Root, cfg.Config}

	fundH := &tx.ChaincodeTx{ccname, TransferHandler(cfg), nil, nil}
	fundH.PreHandlers = append(fundH.PreHandlers, tx.InnerAddrVerifier{ib, FundAddrCred(fundH.Handler.Msg())})
	ret[Method_Transfer] = fundH

	touchH := &tx.ChaincodeTx{ccname, TouchHandler(), nil, nil}
	touchH.PreHandlers = append(touchH.PreHandlers, tx.InnerAddrRegister{ib, FundAddrCred(touchH.Handler.Msg())})
	ret[Method_TouchAddr] = touchH

	return
}

func GeneralInvokingTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	h := TransferHandler(cfg)
	txh := &tx.ChaincodeTx{ccname, h, nil, nil}
	txh.PreHandlers = append(txh.PreHandlers, tx.AddrCredVerifier{FundAddrCred(h.Msg()), nil})
	ret[Method_Transfer] = txh
	return
}

func LimitedQueryTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	ret[Method_QueryToken] = &tx.ChaincodeTx{ccname, TokenQueryHandler(cfg), nil, nil}
	ret[Method_QueryGlobal] = &tx.ChaincodeTx{ccname, GlobalQueryHandler(cfg), nil, nil}
	return
}

func GeneralQueryTemplate(ccname string, cfg LocalConfig) (ret tx.CollectiveTxs) {

	ret = LimitedQueryTemplate(ccname, cfg)
	ret[nonce.Method_Query] = &tx.ChaincodeTx{ccname, nonce.NonceQueryHandler(cfg.Nonce()), nil, nil}

	return
}
