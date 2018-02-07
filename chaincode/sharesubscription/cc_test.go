package subscription

import (
	_ "bytes"
	"github.com/abchain/fabric/core/chaincode/shim"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/crypto"
	tx "hyperledger.abchain.org/tx"
	"math/big"
	"testing"
)

const (
	test_tag    = "test"
	test_ccname = "testCC"

	totalToken = "10000000000000000000000000000"
	assign1    = "500000000000000000000000000"
	assign2    = "800000000000000000000000000"
	assign3    = "200000000000000000000000000"
	rest       = "8500000000000000000000000000"
	trans1     = "400000000000000000000000000" //addr1 -> addr4
	trans2     = "200000000000000000000000000" //addr2 -> addr4
	trans3     = "200000000000000000000000000" //addr3 -> addr1
	transfail  = "300000000000000000000000001" //addr1 -> (fail)
	result1    = "300000000000000000000000000" //addr1
	result2    = "600000000000000000000000000" //addr4
	addr1      = "1984"
	addr2      = "AnimalFarm"
	addr3      = "GeorgeOrwell"
	addr4      = "BigBrother"
)

var spout *GeneralCall
var tokenSpout *token.GeneralCall
var bolt *rpc.DummyCallerBuilder
var stub *shim.MockStub
var tokencfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}
var tokenQuerycfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, true}}
var cfg = &StandardContractConfig{test_tag, false, tokencfg}
var querycfg = &StandardContractConfig{test_tag, true, tokenQuerycfg}

var contract map[string]uint32
var addr1S, addr2S, addr3S, addr4S string

func initContract(t *testing.T) {
	addr1S = tx.NewAddressFromHash([]byte(addr1)).ToString()
	addr2S = tx.NewAddressFromHash([]byte(addr2)).ToString()
	addr3S = tx.NewAddressFromHash([]byte(addr3)).ToString()
	addr4S = tx.NewAddressFromHash([]byte(addr4)).ToString()

	contract = map[string]uint32{
		addr1S: 34,
		addr2S: 35,
		addr3S: 31,
		addr4S: 1,
	}
}

func initTest(t *testing.T) {

	//we only use the db in stub and never do mocking from chaincode interface
	stub = shim.NewMockStub("ShareTest", nil)

	tokenSpout = &token.GeneralCall{txgen.SimpleTxGen(test_ccname)}
	bolt = &rpc.DummyCallerBuilder{test_ccname, stub}

	stub.MockTransactionStart("deployment")

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	deployargs, err := token.CCDeploy(total, nil)

	if err != nil {
		t.Fatal(err)
	}

	var h token.CCDeployHandler = test_tag
	handlers := make(map[string]rpc.DeployHandler)
	handlers[token.DeployMethod] = h

	err = rpc.DeployCC(stub, deployargs, handlers)

	if err != nil {
		t.Fatal(err)
	}

	spout = &GeneralCall{txgen.SimpleTxGen(test_ccname)}
}

func TestContract(t *testing.T) {
	initTest(t)
	initContract(t)

	priv, err := crypto.NewPrivatekey(crypto.DefaultCurveType)
	spout.Credgenerator = txgen.NewSingleKeyCred(priv)

	contractH := NewContractHandler(cfg)
	caller := bolt.GetCaller(contractH)
	bolt.AppendPreHandler(caller, contractH)
	spout.Dispatcher = caller

	stub.MockTransactionStart("contract")

	if err != nil {
		t.Fatal(err)
	}

	addr, err := spout.New(contract, priv.Public())
	if err != nil {
		t.Fatal(err)
	}

	spout.Dispatcher = bolt.GetCaller(QueryHandler(querycfg))

	err, cont := spout.Query(addr)
	if err != nil {
		t.Fatal(err)
	}

	if len(cont.Status) != 4 {
		t.Fatalf("Invalid status count %d", len(cont.Status))
	}

	a1, ok := cont.Status[addr1S]
	if !ok {
		t.Fatal("No record for addr1")
	}

	if a1.Weight < 336633 || a1.Weight > 336634 {
		t.Fatalf("Invalid weight for a1: %d", a1.Weight)
	}

	a4, ok := cont.Status[addr4S]
	if !ok {
		t.Fatal("No record for addr4")
	}

	if a4.Weight < 9900 || a4.Weight > 9901 {
		t.Fatalf("Invalid weight for a4: %d", a4.Weight)
	}

	assignt1, ok := big.NewInt(0).SetString(assign1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	stub.MockTransactionStart("assign")
	tokenSpout.Dispatcher = bolt.GetCaller(token.AssignHandler(tokencfg))

	_, err = tokenSpout.Assign(addr, assignt1)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("redeem1")
	spout.Dispatcher = bolt.GetCaller(RedeemHandler(cfg))

	_, err = spout.Redeem(addr, []byte(addr1), big.NewInt(0))

	tokenSpout.Dispatcher = bolt.GetCaller(token.TokenQueryHandler(tokenQuerycfg))

	err, data1 := tokenSpout.Account([]byte(addr1))
	if err != nil {
		t.Fatal(err)
	}

	dlim, ok := big.NewInt(0).SetString("168316500000000000000000000", 0)
	if !ok {
		t.Fatal("parse int fail")
	}

	ulim, ok := big.NewInt(0).SetString("168317000000000000000000000", 0)
	if !ok {
		t.Fatal("parse int fail")
	}

	bal := big.NewInt(0).SetBytes(data1.Balance)

	if bal.Cmp(dlim) < 0 || bal.Cmp(ulim) > 0 {
		t.Fatalf("wrong redeem amount: %s", bal.String())
	}

	_, err = spout.Redeem(addr, []byte(addr2), bal)

	err, data2 := tokenSpout.Account([]byte(addr2))
	if err != nil {
		t.Fatal(err)
	}

	if bal.Cmp(big.NewInt(0).SetBytes(data2.Balance)) != 0 {
		t.Fatalf("wrong redeem amount for addr2")
	}
}
