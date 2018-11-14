package registrar

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralInvokingTemplate(ccname string, cfg RegistrarConfig) (ret tx.CollectiveTxs) {
	ret = tx.NewCollectiveTxs()

	ret[Method_AdminRegistrar] = &tx.ChaincodeTx{ccname, AdminRegistrarHandler(cfg), nil, nil}
	ret[Method_Registrar] = &tx.ChaincodeTx{ccname, QueryPkHandler(cfg), nil, nil}

	return
}
