package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/impl"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"reflect"
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

//it is safe to return nil for a mal-formed message
//and MUST return empty array if some expected addresses
//is not available (that is, even there is still some address can be
//returned, they MUST NOT shown)
type ListAddresses func(proto.Message) []*txutil.Address

type MsgAddresses interface {
	GetAddresses() []*txutil.Address
}

type AddrVerifier interface {
	TxPreHandler
	Verify(*txutil.Address) error
}

type ClonableAddrVerifier interface {
	Clone() AddrVerifier
}

type AddrCredInspector interface {
	AddVerifier(AddrVerifier)
}

//Verify addresses which the interface required
//One of the interface is used, And interface is tried from top to bottom
type addrCredVerifier struct {
	ListAddresses
	inspectors []AddrVerifier
}

func NewAddrCredVerifier(la ListAddresses) *addrCredVerifier {
	return &addrCredVerifier{la, nil}
}

func AttachAddrVerifier(phs []TxPreHandler, v AddrVerifier) {
	for _, ph := range phs {
		if phA, ok := ph.(AddrCredInspector); ok {
			phA.AddVerifier(v)
		}
	}
}

func cloneAddrVerifier(prototype AddrVerifier) AddrVerifier {
	if cl, ok := prototype.(ClonableAddrVerifier); ok {
		return cl.Clone()
	}

	//the under-lying object of verifier is no way to be immutable
	//(two method has to be called in sequence), so we has to prepare
	//new instance for each calling, if there is not a clone interface,
	//we use reflection to make new instance of the incoming prototype
	return reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface().(AddrVerifier)
}

type failVerifier struct {
	e error
}

func (failVerifier) PreHandling(shim.ChaincodeStubInterface, string, txutil.Parser) error {
	return nil
}

func (v failVerifier) Verify(*txutil.Address) error { return v.e }

func (v *addrCredVerifier) AddVerifier(vv AddrVerifier) { v.inspectors = append(v.inspectors, vv) }

func (v *addrCredVerifier) PreHandling(stub shim.ChaincodeStubInterface, method string, tx txutil.Parser) error {

	var addrs []*txutil.Address

	if v.ListAddresses != nil {
		addrs = v.ListAddresses(tx.GetMessage())
	} else if maddr, ok := tx.GetMessage().(MsgAddresses); ok {
		addrs = maddr.GetAddresses()
	}

	if len(addrs) == 0 {
		return fmt.Errorf("No address is available")
	}

	cred := tx.GetAddrCredential()
	inps := make([]AddrVerifier, len(v.inspectors))

	for _, addr := range addrs {

		done := false
		//verified by inspectors, finnaly by incomming credential
		for i, inp := range inps {
			if inp == nil {
				inp = cloneAddrVerifier(v.inspectors[i])
				if err := inp.PreHandling(stub, method, tx); err != nil {
					inp = failVerifier{err}
				}

				inps[i] = inp
			}

			if inp.Verify(addr) == nil {
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
