package ccprotos

import (
	tx "hyperledger.abchain.org/core/tx"
	"math/big"
)

type simpleFundDetail struct {
	Amount string `json:"amount,omitempty"`
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
}

func (m *SimpleFund) MsgDetail() interface{} {

	ret := new(simpleFundDetail)

	addr, err := tx.NewAddressFromPBMessage(m.From)
	if err == nil {
		ret.From = addr.ToString()
	}

	addr, err = tx.NewAddressFromPBMessage(m.To)
	if err == nil {
		ret.To = addr.ToString()
	}

	ret.Amount = big.NewInt(0).SetBytes(m.Amount).String()

	return ret
}
