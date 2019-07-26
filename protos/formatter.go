package protos

import (
	"encoding/json"
)

var TxAddrMarshaller = func(m *TxAddr) ([]byte, error) {
	return json.Marshal(m.Hash)
}

func (m *TxAddr) MarshalJSON() ([]byte, error) {
	return TxAddrMarshaller(m)
}
