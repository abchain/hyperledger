package ccprotos

import (
	tx "hyperledger.abchain.org/core/tx"
)

func (m *RedeemContract) GetAddresses() (addrs []*tx.Address) {
	for _, redeemAddr := range m.GetRedeems() {
		if addr, err := tx.NewAddressFromPBMessage(redeemAddr); err == nil {
			addrs = append(addrs, addr)
		}
	}

	return
}

func (m *RegContract) GetAddresses() (addrs []*tx.Address) {

	if addr, err := tx.NewAddressFromPBMessage(m.GetDelegator()); err == nil {
		addrs = append(addrs, addr)
	}

	return
}
