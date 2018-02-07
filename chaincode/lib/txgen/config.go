package tx

import (
	_ "errors"
	"hyperledger.abchain.org/crypto"
)

func SimpleTxGen(ccname string) *TxGenerator {

	return &TxGenerator{nil, nil, nil, nil, nil, ccname}
}

func DefaultTxGen(ccname string, privkey *crypto.PrivateKey) *TxGenerator {

	return &TxGenerator{nil, nil, NewSingleKeyCred(privkey), nil, nil, ccname}
}
