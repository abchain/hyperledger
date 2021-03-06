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

	verifier := tx.NewAddrCredVerifier(nil)

	txh := &tx.ChaincodeTx{ccname, TransferHandler(cfg), nil, nil}
	txh.PreHandlers = append(txh.PreHandlers, verifier)
	ret[Method_Transfer] = txh

	txh = &tx.ChaincodeTx{ccname, TransferHandler2(cfg), nil, nil}
	txh.PreHandlers = append(txh.PreHandlers, verifier)
	ret[Method_Transfer2] = txh

	return
}

func ExtendInvokingTemplate(cts tx.CollectiveTxs, verifier tx.AddrVerifier) tx.CollectiveTxs {

	if h, ok := cts[Method_Transfer]; ok {
		tx.AttachAddrVerifier(h.PreHandlers, verifier)
	}

	if h, ok := cts[Method_Transfer2]; ok {
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

//the local config must provide both a executable interface and the sub-config (corresponding its sub interface)
//for local handler building
type LocalConfig interface {
	TokenConfig
	Nonce() nonce.NonceConfig
}

func ExtendedQueryTemplate(ccname string, cfg LocalConfig) (ret tx.CollectiveTxs) {

	ret = GeneralQueryTemplate(ccname, cfg)
	ret[nonce.Method_Query] = &tx.ChaincodeTx{ccname, nonce.NonceQueryHandler(cfg.Nonce()), nil, nil}

	return
}
