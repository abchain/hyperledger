package addrspace

import (
	"hyperledger.abchain.org/chaincode/lib/txhandle"
)

func GeneralTemplate(ccname string, cfg AddrSpaceConfig) (ret tx.CollectiveTxs) {

	ret = tx.NewCollectiveTxs()
	ret[Method_Reg] = &tx.ChaincodeTx{ccname, RegHandler(cfg), nil, nil}
	ret[Method_Query] = &tx.ChaincodeTx{ccname, QueryHandler(cfg), nil, nil}

	return
}

func DummyTemplate(ccname string) (ret tx.CollectiveTxs) {
	return GeneralTemplate(ccname, DummyImplCfg())
}
