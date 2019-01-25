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

	h := TransferHandler(cfg)
	txh := &tx.ChaincodeTx{ccname, h, nil, nil}
	txh.PreHandlers = append(txh.PreHandlers, tx.NewAddrCredVerifier(nil))
	ret[Method_Transfer] = txh
	return
}

func ExtendInvokingTemplate(cts tx.CollectiveTxs, ccname string, cfg *StandardTokenConfig) tx.CollectiveTxs {

	ib := &tx.InnerAddrBase{Root: cfg.Root, Config: cfg.Config}

	if h, ok := cts[Method_Transfer]; ok {
		tx.AttachAddrVerifier(h.PreHandlers, &tx.InnerAddrVerifier{InnerAddrBase: ib})
	}

	touchH := &tx.ChaincodeTx{ccname, TouchHandler(cfg), nil, nil}
	touchH.PostHandlers = append(touchH.PostHandlers, tx.InnerAddrRegister{ib, nil})
	cts[Method_TouchAddr] = touchH

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
