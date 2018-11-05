package tx

import (
	_ "errors"
	"hyperledger.abchain.org/core/crypto"
)

func SimpleTxGen(ccname string) *TxGenerator {

	return &TxGenerator{Ccname: ccname}
}

func DefaultTxGen(ccname string, privkey *crypto.PrivateKey) *TxGenerator {

	return &TxGenerator{Credgenerator: NewSingleKeyCred(privkey), Ccname: ccname}
}
