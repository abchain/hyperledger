package ccprotos

import (
	"hyperledger.abchain.org/core/crypto"
	tx "hyperledger.abchain.org/core/tx"
)

func (m *RegPublicKey) GetAddresses() []*tx.Address {

	pk, err := crypto.PublicKeyFromBytes(m.PkBytes)
	if err != nil {
		return nil
	}

	addr, err := tx.NewAddress(pk)

	if err != nil {
		return nil
	}
	return []*tx.Address{addr}
}
