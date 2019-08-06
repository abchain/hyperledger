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
type ListAddresses func(shim.ChaincodeStubInterface, proto.Message) []*txutil.Address

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
	AddVerifier(AddrVerifier) bool
}

var Verifier_CheckCred = fmt.Errorf("Pass if credential can be verified")

//Verify addresses which the interface required
//One of the interface is used, And interface is tried from top to bottom
type addrCredVerifier struct {
	ListAddresses
	inspectors []AddrVerifier
}

func EmptyAddrCredVerifier(la ListAddresses) *addrCredVerifier {
	return &addrCredVerifier{la, nil}
}

//create default verifier
func NewAddrCredVerifier(la ListAddresses) *addrCredVerifier {
	return &addrCredVerifier{la, []AddrVerifier{defaultVerifier}}
}

func NewAddrCredVerifierFromTemplate(la ListAddresses, tp TxPreHandler) *addrCredVerifier {

	tpV := tp.(*addrCredVerifier)
	inspCpy := make([]AddrVerifier, len(tpV.inspectors))
	copy(inspCpy, tpV.inspectors)

	return &addrCredVerifier{la, inspCpy}
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
	//we use reflection to make new instance of the incoming prototype,
	//the prototype is assigned to the created one
	protoV := reflect.Indirect(reflect.ValueOf(prototype))
	newV := reflect.New(protoV.Type())
	reflect.Indirect(newV).Set(protoV)

	return newV.Interface().(AddrVerifier)
}

type failVerifier struct {
	e error
}

func (failVerifier) PreHandling(shim.ChaincodeStubInterface, string, txutil.Parser) error {
	return nil
}

func (v failVerifier) Verify(*txutil.Address) error { return v.e }

func (v *addrCredVerifier) AddDefaultVerifier() {
	v.AddVerifier(defaultVerifier)
}

func (v *addrCredVerifier) AddVerifier(vv AddrVerifier) bool {

	//verify ...
	for _, inp := range v.inspectors {
		if inp == vv {
			return false
		}
	}

	v.inspectors = append(v.inspectors, vv)
	return true
}

var defaultVerifier = failVerifier{Verifier_CheckCred}
var noCredential = fmt.Errorf("No credential inspector")

func (v *addrCredVerifier) verifyCore(stub shim.ChaincodeStubInterface,
	method string, tx txutil.Parser,
	addrs []*txutil.Address, inps []AddrVerifier) (error, []AddrVerifier) {

	cred := tx.GetAddrCredential()

	for _, addr := range addrs {

		verifyErr := noCredential
		//verified by inspectors, finnaly by incomming credential

		traceLog := []string{}
		for i, inp := range inps {
			if inp == nil {
				inp = cloneAddrVerifier(v.inspectors[i])
				if err := inp.PreHandling(stub, method, tx); err != nil {
					inp = failVerifier{err}
				}

				inps[i] = inp
			}

			verifyErr = inp.Verify(addr)
			if verifyErr == Verifier_CheckCred && cred != nil {
				verifyErr = cred.Verify(*addr)
			}

			if verifyErr == nil {
				break
			} else {
				traceLog = append(traceLog, verifyErr.Error()+"\n")
			}
		}

		if verifyErr != nil {
			return fmt.Errorf("addr [%s] verify failure: \n%v", addr.ToString(), traceLog), inps
		}
	}

	return nil, inps

}

func (v *addrCredVerifier) PreHandling(stub shim.ChaincodeStubInterface, method string, tx txutil.Parser) error {

	var addrs []*txutil.Address

	if v.ListAddresses != nil {
		addrs = v.ListAddresses(stub, tx.GetMessage())
	} else if maddr, ok := tx.GetMessage().(MsgAddresses); ok {
		addrs = maddr.GetAddresses()
	}

	if len(addrs) == 0 {
		return fmt.Errorf("No address is available")
	}

	err, _ := v.verifyCore(stub, method, tx, addrs, make([]AddrVerifier, len(v.inspectors)))

	return err
}

type addrOrCredVerifier struct {
	*addrCredVerifier
}

func NewAddrOrCredVerifier(la ListAddresses) addrOrCredVerifier {
	return addrOrCredVerifier{&addrCredVerifier{la, []AddrVerifier{defaultVerifier}}}
}

func (v addrOrCredVerifier) PreHandling(stub shim.ChaincodeStubInterface, method string, tx txutil.Parser) error {

	var addrs []*txutil.Address

	if v.ListAddresses != nil {
		addrs = v.ListAddresses(stub, tx.GetMessage())
	} else if maddr, ok := tx.GetMessage().(MsgAddresses); ok {
		addrs = maddr.GetAddresses()
	}

	if len(addrs) == 0 {
		return fmt.Errorf("No address is available")
	}

	var allErrors []error
	inps := make([]AddrVerifier, len(v.inspectors))
	for _, addr := range addrs {
		var err error
		err, inps = v.verifyCore(stub, method, tx, []*txutil.Address{addr}, inps)
		if err == nil {
			return err
		} else {
			allErrors = append(allErrors, err)
		}
	}

	return fmt.Errorf("All address verify fail: %v", allErrors)
}
