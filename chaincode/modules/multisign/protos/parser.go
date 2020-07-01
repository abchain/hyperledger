package ccauthprotos

import (
	"encoding/json"
	"fmt"
	txutil "hyperledger.abchain.org/core/tx"
)

//Marshaller of JSON for parsing
func (c ctaddr) MarshalJSON() ([]byte, error) {

	if c.Addr == nil {
		return json.Marshal("Wrong contract member")
	}

	caddr := txutil.NewAddressFromHash(c.Addr)

	return json.Marshal(fmt.Sprintf("%s:%d", caddr.ToString(), c.Weight))

}

//GetAddresses indicate we use current contract address
func (m *Update) GetAddresses() []*txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(m.Addr)
	if err != nil {
		return nil
	}

	return []*txutil.Address{addr}
}
