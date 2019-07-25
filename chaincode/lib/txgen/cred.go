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

	sig, err := c.privkey.Sign(builder.GetHash())
	if err != nil {
		return err
	}

	builder.GetCredBuilder().AddSignature(sig)

	return nil
}

func NewSingleKeyCred(privkey crypto.Signer) TxCredHandler {
	return &singleKeyCred{privkey}
}

type multiKeyCred struct {
	// privkey *crypto.PrivateKey

	privkeys []crypto.Signer
}

func (c *multiKeyCred) DoCred(builder txutil.Builder) error {

	for _, key := range c.privkeys {
		sig, err := key.Sign(builder.GetHash())
		if err != nil {
			return err
		}

		builder.GetCredBuilder().AddSignature(sig)

	}

	return nil
}

func NewMultiKeyCred(privkey ...crypto.Signer) TxCredHandler {
	return &multiKeyCred{privkey}
}
