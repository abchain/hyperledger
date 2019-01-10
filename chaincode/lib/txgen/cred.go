package tx

import (
	_ "errors"

	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

type TxCredHandler interface {
	DoCred(txutil.Builder) error
}

type singleKeyCred struct {
	// privkey *crypto.PrivateKey

	privkey crypto.Signer
}

func (c *singleKeyCred) DoCred(builder txutil.Builder) error {
	return builder.Sign(c.privkey)
}

func NewSingleKeyCred(privkey crypto.Signer) TxCredHandler {
	return &singleKeyCred{privkey}
}
