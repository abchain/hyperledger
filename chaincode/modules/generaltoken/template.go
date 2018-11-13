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

func GeneralInvokingTemplate(ccname string, cfg TokenConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	fundTx := &tx.ChaincodeTx{ccname, TransferHandler(cfg), nil, nil}
	//only append address credverify (tx must signed by the to address)
	fundTx.PreHandlers = append(fundTx.PreHandlers, tx.AddrCredVerifier{FundAddrCred(fundTx.Handler.Msg()), nil})
	ret[Method_Transfer] = fundTx
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
