package subscription

import (
	_ "bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/core/crypto/ecdsa"
	tx "hyperledger.abchain.org/core/tx"
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

var bolt *rpc.ChaincodeAdapter
var tokenbolt *rpc.ChaincodeAdapter

var contract map[string]int32
var addr1S, addr2S, addr3S, addr4S string

func init() {
	addr1S = tx.NewAddressFromHash([]byte(addr1)).ToString()
	addr2S = tx.NewAddressFromHash([]byte(addr2)).ToString()
	addr3S = tx.NewAddressFromHash([]byte(addr3)).ToString()
	addr4S = tx.NewAddressFromHash([]byte(addr4)).ToString()

	contract = map[string]int32{
		addr1S: 34,
		addr2S: 35,
		addr3S: 31,
		addr4S: 1,
	}
}

func initCond(mutilcc bool) {

	cfg := NewConfig(test_tag)
	querycfg := NewConfig(test_tag)
	querycfg.SetReadOnly(true)

	tokencfg := cfg.TokenCfg
	tokenQuerycfg := querycfg.TokenCfg

	if mutilcc {
		tokencfg = token.NewConfig(test_tag)
		tqc := token.NewConfig(test_tag)
		tqc.SetReadOnly(true)
		tokenQuerycfg = tqc
	}

	shareCC := GeneralInvokingTemplate(test_ccname, cfg).MustMerge(GeneralQueryTemplate(test_ccname, querycfg))

	tokenCC := token.GeneralAdminTemplate(test_ccname, tokencfg)
	tokenCC = tokenCC.MustMerge(token.GeneralQueryTemplate(test_ccname, tokenQuerycfg))

	if mutilcc {

		tokenCfg := token.NewConfig(test_tag)
		tokenCC = tokenCC.MustMerge(token.GeneralInvokingTemplate(test_ccname, tokenCfg))
		tokenCC = token.ExtendInvokingTemplate(tokenCC, test_ccname, tokenCfg)

		tokenbolt = rpc.NewLocalChaincode(txhandle.CollectiveTxs_InnerSupport(tokenCC))
		tokenbolt.Name = "tokenCC"
		cfg.TokenCfg = token.InnerInvokeConfig{txgen.InnerChaincode(tokenbolt.Name)}
		querycfg.TokenCfg = cfg.TokenCfg

		bolt = rpc.NewLocalChaincode(shareCC)
		bolt.Name = "contractCC"
		bolt.Invokables[tokenbolt.Name] = tokenbolt.MockStub

	} else {

		bolt = rpc.NewLocalChaincode(shareCC.MustMerge(tokenCC))
		tokenbolt = bolt
	}

}

func initTest(t *testing.T) {

	spoutcore := txgen.SimpleTxGen(test_ccname)
	tokenSpout := &token.GeneralCall{spoutcore}
	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginDeploy(nil)

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	err := tokenSpout.Init(total)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	initCond(false)
	initTest(t)
}

func TestInitMutliCC(t *testing.T) {
	initCond(true)
	initTest(t)
}

func testContractBase(t *testing.T) {

	spoutcore := txgen.SimpleTxGen(test_ccname)
	tokenspoutcore := txgen.SimpleTxGen(test_ccname)
	spoutcore.Dispatcher = bolt
	tokenspoutcore.Dispatcher = tokenbolt

	tokenSpout := &token.GeneralCall{tokenspoutcore}
	spout := &GeneralCall{spoutcore}

	priv, err := ecdsa.NewDefaultPrivatekey()
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(priv)

	if err != nil {
		t.Fatal(err)
	}

	privAddr, err := tx.NewAddress(priv.Public())
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.BeginTx(nil)
	bolt.SpecifyTxID("contract")

	addr, err := spout.New(contract, privAddr.Hash)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	err, cont := spout.Query(addr)
	if err != nil {
		t.Fatal(err)
	}

	if len(cont.Status) != 4 {
		t.Fatalf("Invalid status count %d", len(cont.Status))
	}

	a1, ok := cont.Find(addr1S)
	if !ok {
		t.Fatal("No record for addr1")
	}

	if a1.Weight < 336633 || a1.Weight > 336634 {
		t.Fatalf("Invalid weight for a1: %d", a1.Weight)
	}

	a4, ok := cont.Find(addr4S)
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

	tokenspoutcore.BeginTx(nil)
	bolt.SpecifyTxID("assign")

	_, err = tokenSpout.Assign(addr, assignt1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tokenspoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.BeginTx(nil)
	bolt.SpecifyTxID("redeem1")

	_, err = spout.Redeem(addr, big.NewInt(0), [][]byte{[]byte(addr1)})
	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

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

	bal := data1.Balance

	if bal.Cmp(dlim) < 0 || bal.Cmp(ulim) > 0 {
		t.Fatalf("wrong redeem amount: %s", bal.String())
	}

	spoutcore.BeginTx(nil)
	bolt.SpecifyTxID("redeem2")
	_, err = spout.Redeem(addr, bal, [][]byte{[]byte(addr2)})
	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	err, data2 := tokenSpout.Account([]byte(addr2))
	if err != nil {
		t.Fatal(err)
	}

	if bal.Cmp(data2.Balance) != 0 {
		t.Fatalf("wrong redeem amount for addr2")
	}
}

func TestContract(t *testing.T) {
	TestInit(t)
	testContractBase(t)
}

func TestContract_Multicc(t *testing.T) {
	TestInitMutliCC(t)
	testContractBase(t)
}
