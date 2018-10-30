package tx

import (
	"fmt"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/tx"
	"strings"
)

type TxAttrVerifier map[string]string
type TxMultiAttrVerifier map[string][]string

func (req TxAttrVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, _ txutil.Parser) error {

	for attrkey, expect := range req {
		attr, err := stub.GetCallerAttribute(attrkey)
		if err != nil {
			return err
		}

		if strings.Compare(string(attr), expect) != 0 {
			return fmt.Errorf("No rivilege for attr %s, provide:[%s]", attrkey, attr)
		}
	}

	return nil
}

func (req TxMultiAttrVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, _ txutil.Parser) error {

	for attrkey, expects := range req {
		attr, err := stub.GetCallerAttribute(attrkey)
		if err != nil {
			return err
		}
		matched := false

		for _, expect := range expects {
			if strings.Compare(string(attr), expect) == 0 {
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

//Verify only one address and its corresponding cred
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

	var e error

	if v.ParseAddress != nil {
		if err := tryAddrParser(v.ParseAddress, cred); err == nil {
			return nil
		} else {
			e = err
		}
	}

	if v.MatchAddress != nil {
		if err := tryAddrMatcher(v.MatchAddress, cred); err == nil {
			return nil
		} else {
			e = err
		}
	}

	return e
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
		if err == nil && v.Match(addr) && cred.Verify(*addr) == nil {
			//match!
			return nil
		}
	}

	return fmt.Errorf("No valid creds found")
}
