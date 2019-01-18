package tx

import (
	"fmt"
	"hyperledger.abchain.org/chaincode/impl"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"
)

type TxAttrVerifier map[string]string
type TxMultiAttrVerifier map[string][]string

func (req TxAttrVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, _ txutil.Parser) error {

	attrif, err := impl.GetCallerAttributes(stub)
	if err != nil {
		return err
	}

	for attrkey, expect := range req {
		attr, err := attrif.GetCallerAttribute(attrkey)
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

	attrif, err := impl.GetCallerAttributes(stub)
	if err != nil {
		return err
	}

	for attrkey, expects := range req {
		attr, err := attrif.GetCallerAttribute(attrkey)
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

type ListAddresses interface {
	ListAddress() []*txutil.Address
}

type AddrVerifier interface {
	Verify(*txutil.Address) error
}

type AddrCredInspector interface {
	AddVerifier(AddrVerifier)
}

//Verify addresses which the interface required
//One of the interface is used, And interface is tried from top to bottom
type addrCredVerifier struct {
	ParseAddress
	ListAddresses
	inspectors []AddrVerifier
}

func NewAddrCredVerifier(pa ParseAddress, la ListAddresses) *addrCredVerifier {
	return &addrCredVerifier{pa, la, nil}
}

func AttachAddrVerifier(phs []TxPreHandler, v AddrVerifier) {
	for _, ph := range phs {
		if phA, ok := ph.(AddrCredInspector); ok {
			phA.AddVerifier(v)
		}
	}
}

func (v *addrCredVerifier) AddVerifier(vv AddrVerifier) { v.inspectors = append(v.inspectors, vv) }

func (v *addrCredVerifier) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	if v.ParseAddress == nil && v.ListAddresses == nil {
		panic("Uninit interface")
	}

	var addrs []*txutil.Address
	if v.ParseAddress != nil {
		addrs = append(addrs, v.GetAddress())
	} else {
		addrs = v.ListAddress()
	}

	cred := tx.GetAddrCredential()

	for _, addr := range addrs {

		done := false
		//verified by inspectors, finnaly by incomming credential
		for _, inp := range v.inspectors {
			if inp.Verify(addr) != nil {
				done = true
				break
			}
		}

		if !done && cred != nil {
			if cred.Verify(*addr) == nil {
				done = true
			}
		}

		if !done {
			return fmt.Errorf("Addr [%s] has no credential", addr.ToString())
		}
	}

	return nil
}
