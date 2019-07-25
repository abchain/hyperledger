package addrspace

import (
	"bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
)

const (
	test_tag     = "test"
	test_ccname  = "testCC"
	workerCCName = "worker"
)

type dynamicCC func(shim.ChaincodeStubInterface)

func (f dynamicCC) Invoke(stub shim.ChaincodeStubInterface,
	_ string, _ [][]byte, _ bool) ([]byte, error) {
	f(stub)
	return []byte{}, nil
}

var bolt *rpc.ChaincodeAdapter

func initCond() {
	invokeCfg := NewConfig(test_tag)
	cc := GeneralTemplate(test_ccname, invokeCfg)

	bolt = rpc.NewLocalChaincode(txhandle.CollectiveTxs_InnerSupport(cc))
	bolt.Name = workerCCName
}

func ccDynamicCall(ccname string, cc dynamicCC, t *testing.T) {

	callcc := rpc.NewLocalChaincode(cc)
	callcc.Name = ccname
	callcc.Invokables[bolt.Name] = bolt.MockStub

	caller := txgen.TxGenerator{Ccname: test_ccname, Dispatcher: callcc}
	caller.Ccname = test_ccname
	caller.BeginTx(nil)
	defer caller.Result()
	err := caller.Invoke("Any", new(empty.Empty))
	if err != nil {
		t.Fatal(err)
	}

}

func TestReg(t *testing.T) {

	initCond()

	spoutCfg := InnerinvokeImpl(txgen.InnerChaincode(workerCCName))

	var err error
	regCC := func(stub shim.ChaincodeStubInterface) {
		spout := spoutCfg.NewTx(stub, []byte{0})
		err = spout.RegisterCC()
	}

	ccDynamicCall("cc1", regCC, t)
	if err != nil {
		t.Fatal(t)
	}

	var prefix []byte

	queryCC := func(stub shim.ChaincodeStubInterface) {
		spout := spoutCfg.NewTx(stub, []byte{0})

		prefix, err = spout.QueryPrefix()
		if err != nil {
			t.Fatal(t)
		}
	}

	ccDynamicCall("cc1", queryCC, t)

	if len(prefix) == 0 {
		t.Fatal("empty prefix")
	}
	p1 := prefix

	ccDynamicCall("cc1", queryCC, t)
	p2 := prefix

	if bytes.Compare(p1, p2) != 0 {
		t.Fatalf("different prefix result: %x vs %x", p1, p2)
	}

	ccDynamicCall("cc1", regCC, t)
	if err == nil {
		t.Fatal("can duplicated")
	}

	ccDynamicCall("cc2", queryCC, t)
	if bytes.Compare(p1, prefix) != 0 {
		t.Fatalf("[cache not work] different prefix result: %x", prefix)
	}

	spoutCfg = InnerinvokeImpl(txgen.InnerChaincode(workerCCName))

	ccDynamicCall("cc2", regCC, t)
	if err != nil {
		t.Fatal(err)
	}

	ccDynamicCall("cc2", queryCC, t)
	p3 := prefix
	if len(p3) == 0 {
		t.Fatal("empty prefix")
	}

	if bytes.Compare(p1, p3) == 0 {
		t.Fatal("same prefix for different cc")
	}

}

func TestNormalizeAddr(t *testing.T) {

	initCond()

	spoutCfg := InnerinvokeImpl(txgen.InnerChaincode(workerCCName))

	var err error

	regCC := func(stub shim.ChaincodeStubInterface) {
		spout := spoutCfg.NewTx(stub, []byte{0})
		err = spout.RegisterCC()
	}

	ccDynamicCall("cc1", regCC, t)
	if err != nil {
		t.Fatal(t)
	}

	ccDynamicCall("cc2", regCC, t)
	if err != nil {
		t.Fatal(t)
	}

	var addr []byte
	normalizeCC := func(stub shim.ChaincodeStubInterface) {
		spout := spoutCfg.NewTx(stub, []byte{0})
		addr, err = spout.NormalizeAddress(addr)
	}

	addr = []byte{42, 42}

	ccDynamicCall("cc1", normalizeCC, t)
	if err != nil {
		t.Fatal(t)
	}

	t.Logf("[%X]", addr)

	if len(addr) <= 2 {
		t.Fatal("Addr is not normalized")
	}

	addr1 := addr
	addr = []byte{42, 42}
	spoutCfg = InnerinvokeImpl(txgen.InnerChaincode(workerCCName))

	ccDynamicCall("cc2", normalizeCC, t)
	if err != nil {
		t.Fatal(err)
	}

	if len(addr) <= 2 {
		t.Fatal("Addr is not normalized")
	}

	t.Logf("[%X]", addr)

	if bytes.Compare(addr, addr1) == 0 {
		t.Fatal("Normalized different cc into same addr")
	}

}
