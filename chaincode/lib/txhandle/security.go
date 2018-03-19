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

type ParseAddress interface {
	GetAddress() *txutil.Address
}

//matchAddress do not require the interface has immutable nature
//side-effect is allowed for "anchored" the matched/unmatched result for the following process
type MatchAddress interface {
	Match(*txutil.Address) bool
}

//One of the interface is used, And interface is tried from top to bottom
type AddrCredVerifier struct {
	ParseAddress
	MatchAddress
}

func (v AddrCredVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	if v.ParseAddress == nil && v.MatchAddress == nil {
		panic("Uninit interface")
	}

	cred := tx.GetAddrCredential()

	if cred == nil {
		return fmt.Errorf("Tx contains no credentials")
	}

	if tryAddrParser(v.ParseAddress, cred) == nil {
		return nil
	}

	if err := tryAddrMatcher(v.MatchAddress, cred); err != nil {
		return err
	}

	return nil
}

func tryAddrParser(v ParseAddress, cred txutil.AddrCredentials) error {

	addr := v.GetAddress()

	if addr == nil {
		return fmt.Errorf("Invalid address")
	}

	return cred.Verify(*addr)
}

func tryAddrMatcher(v MatchAddress, cred txutil.AddrCredentials) error {

	allpks := cred.ListCredPubkeys()

	for _, pk := range allpks {
		addr, err := txutil.NewAddress(pk)
		if err == nil && v.Match(addr) {
			//match!
			return nil
		}
	}

	return fmt.Errorf("No valid creds found")
}
