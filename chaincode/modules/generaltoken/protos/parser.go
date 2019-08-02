package ccprotos

import (
	tx "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/protos"
	"math/big"
)

type baseTokenDetail struct {
	Amount *big.Int `json:"totalTokens,omitempty"`
	Name   string   `json:"token,omitempty"`
}

func (m *BaseToken) MsgDetail() interface{} {

	return &baseTokenDetail{
		Amount: big.NewInt(0).SetBytes(m.GetTotalTokens()),
	}
}

type simpleFundDetail struct {
	Amount *big.Int       `json:"amount,omitempty"`
	From   *protos.TxAddr `json:"from,omitempty"`
	To     *protos.TxAddr `json:"to,omitempty"`
	Name   string         `json:"token,omitempty"`
}

func (m *SimpleFund) MsgDetail() interface{} {

	return &simpleFundDetail{Amount: big.NewInt(0).SetBytes(m.GetAmount()),
		From: m.GetFrom(),
		To:   m.GetTo()}
}

//the default address extractor in message
func (m *SimpleFund) GetAddresses() []*tx.Address {
	addr, err := tx.NewAddressFromPBMessage(m.From)
	if err != nil {
		return nil
	}

	return []*tx.Address{addr}
}

func (m *QueryToken) GetAddresses() []*tx.Address {
	addr, err := tx.NewAddressFromPBMessage(m.Addr)
	if err != nil {
		return nil
	}

	return []*tx.Address{addr}
}

func (m *MultiTokenMsg) GetAddresses() []*tx.Address {

	if m.Msg == nil {
		return nil
	}

	switch em := m.Msg.(type) {
	case *MultiTokenMsg_Fund:
		return em.Fund.GetAddresses()
	case *MultiTokenMsg_Query:
		return em.Query.GetAddresses()
	default:
		return nil
	}

}

func (m *MultiTokenMsg) MsgDetail() interface{} {

	if m.Msg != nil {
		switch em := m.Msg.(type) {
		case *MultiTokenMsg_Fund:
			detail := em.Fund.MsgDetail().(*simpleFundDetail)
			detail.Name = m.GetTokenName()
			return detail
		case *MultiTokenMsg_Init:
			detail := em.Init.MsgDetail().(*baseTokenDetail)
			detail.Name = m.GetTokenName()
			return detail
		}
	}

	return m.GetTokenName() + ": Unknown message"
}
