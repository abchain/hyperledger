package nonce

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralTemplate(ccname string, cfg NonceConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()

	ret[Method_Add] = &tx.ChaincodeTx{ccname, NonceAddHandler(cfg), nil, nil}
	ret[Method_Query] = &tx.ChaincodeTx{ccname, NonceQueryHandler(cfg), nil, nil}

	return
}
