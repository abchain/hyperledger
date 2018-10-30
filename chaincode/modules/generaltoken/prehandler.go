package generaltoken

import (
	txutil "hyperledger.abchain.org/tx"
)

//and set it as RegistrarPreHandler for registrar
func (m *FundMsg) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(m.msg.From)
	if err != nil {
		return nil
	}

	return addr
}
