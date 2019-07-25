package ccprotos

import (
	tx "hyperledger.abchain.org/core/tx"
)

func (m *RegContract) GetAddresses() []*tx.Address {
	addr, err := tx.NewAddressFromPBMessage(m.GetDelegatorAddr())
	if err != nil {
		return nil
	}

	return []*tx.Address{addr}
}

func (m *RedeemContract) GetAddresses() (addrs []*tx.Address) {
	for _, redeemAddr := range m.GetRedeems() {
		if addr, err := tx.NewAddressFromPBMessage(redeemAddr); err == nil {
			addrs = append(addrs, addr)
		}
	}

	return
}
