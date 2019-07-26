package ccprotos

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
