package subscription

// import (
// 	"hyperledger.abchain.org/chaincode/lib/txhandle"
// 	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
// )

// func GeneralInvokingTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

// 	ret = tx.NewCollectiveTxs()

// 	newContractH := share.NewContractHandler(cfg)
// 	newContractTx := &tx.ChaincodeTx{CC_NAME, newContractH, nil, nil}
// 	newContractTx.PreHandlers = append(newContractTx.PreHandlers, newContractH)

// 	fundH := TransferHandler(cfg)
// 	fundTx := &tx.ChaincodeTx{ccname, fundH, nil, nil}
// 	//only append address credverify (tx must signed by the to address)
// 	fundTx.PreHandlers = append(fundTx.PreHandlers, tx.AddrCredVerifier{fundH, nil})
// 	ret[Method_Transfer] = fundTx
// 	return
// }

// func GeneralQueryTemplate(ccname string, cfg ContractConfig) (ret tx.CollectiveTxs) {

// 	ret = tx.NewCollectiveTxs()

// 	ret[Method_QueryToken] = &tx.ChaincodeTx{ccname, TokenQueryHandler(cfg), nil, nil}
// 	ret[Method_QueryGlobal] = &tx.ChaincodeTx{ccname, GlobalQueryHandler(cfg), nil, nil}
// 	ret[nonce.Method_Query] = &tx.ChaincodeTx{ccname, nonce.NonceQueryHandler(cfg.Nonce()), nil, nil}

// 	return
// }
