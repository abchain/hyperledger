package ccshareprotos

import (
	"encoding/json"
	"fmt"
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

func (m *RegContract) GetAddresses() []*tx.Address {

	if addr, err := tx.NewAddressFromPBMessage(m.GetDelegator()); err == nil {
		return []*tx.Address{addr}
	}

	return nil
}

func (m *QueryContract) GetAddresses() []*tx.Address {

	if maddr := m.GetMemberAddr(); maddr != nil {
		if addr, err := tx.NewAddressFromPBMessage(maddr); err == nil {
			return []*tx.Address{addr}
		}
	} else {
		if addr, err := tx.NewAddressFromPBMessage(m.GetContractAddr()); err == nil {
			return []*tx.Address{addr}
		}
	}

	return nil
}

func (m *RedeemResponse) MarshalJSON() ([]byte, error) {

	var ret []string
	for _, nc := range m.GetNonces() {
		ret = append(ret, fmt.Sprintf("%X", nc))
	}

	if len(ret) == 0 {
		return json.Marshal(nil)
	} else {
		return json.Marshal(ret)
	}
}
