package tx

import (
	_ "errors"
	"hyperledger.abchain.org/crypto"
	txutil "hyperledger.abchain.org/tx"
)

type TxCredHandler interface {
	DoCred(txutil.Builder) error
}

type singleKeyCred struct {
	privkey *crypto.PrivateKey
}

func (c *singleKeyCred) DoCred(builder txutil.Builder) error {
	return builder.Sign(c.privkey)
}

func NewSingleKeyCred(privkey *crypto.PrivateKey) TxCredHandler {
	return &singleKeyCred{privkey}
}
