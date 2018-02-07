package tx

import (
	"fmt"
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/lib/util"
	txutil "hyperledger.abchain.org/tx"
	"strings"
)

type TxAttrVerifier map[string]string
type TxMultiAttrVerifier map[string][]string

func (req TxAttrVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, _ txutil.Parser) error {

	for attrkey, expect := range req {
		attr := util.GetAttributes(stub, attrkey)

		if strings.Compare(attr, expect) != 0 {
			return fmt.Errorf("No rivilege for attr %s, provide:[%s]", attrkey, attr)
		}
	}

	return nil
}

func (req TxMultiAttrVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, _ txutil.Parser) error {

	for attrkey, expects := range req {
		attr := util.GetAttributes(stub, attrkey)
		matched := false

		for _, expect := range expects {
			if strings.Compare(attr, expect) == 0 {
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("No rivilege for attr %s, provide:[%s]", attrkey, attr)
		}
	}

	return nil
}
