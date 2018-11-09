package subscription

import (
	_ "bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/core/crypto"
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

var spoutcore = txgen.SimpleTxGen(test_ccname)
var bolt = &rpc.DummyCallerBuilder{CCName: test_ccname}
var tokencfg = token.NewConfig(test_tag)
var tokenQuerycfg = token.NewConfig(test_tag)
var cfg = NewConfig(test_tag)
var querycfg = NewConfig(test_tag)

func init() {
	tokenQuerycfg.SetReadOnly(true)
	querycfg.SetReadOnly(true)
}

var contract map[string]int32
var addr1S, addr2S, addr3S, addr4S string

func initContract(t *testing.T) {
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

func initTest(t *testing.T) {

	//we only use the db in stub and never do mocking from chaincode interface
	bolt.Reset()

	deployTx := &txgen.BatchTxCall{TxGenerator: spoutcore}
	tokenSpout := &token.GeneralCall{deployTx}
	deployTx.BeginDeploy(nil)

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	tokenSpout.Init(total)

	spoutcore.Dispatcher = bolt.GetCaller("deployment",
		txhandle.BatchTxHandler(map[string]*txhandle.ChaincodeTx{
			token.Method_Init: &txhandle.ChaincodeTx{test_ccname, token.InitHandler(tokencfg), nil, nil},
		}))

	err := deployTx.CommitBatch("init")

	if err != nil {
		t.Fatal(err)
	}

	_, err = deployTx.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}
}

func TestContract(t *testing.T) {
	initTest(t)
	initContract(t)
	tokenSpout := &token.GeneralCall{spoutcore}
	spout := &GeneralCall{spoutcore, false}

	priv, err := crypto.NewPrivatekey(crypto.DefaultCurveType)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(priv)

	contractH := NewContractHandler(cfg)
	spoutcore.BeginTx(nil)
	spoutcore.Dispatcher = bolt.GetCaller("contract", contractH)
	bolt.AppendPreHandler(contractH)

	if err != nil {
		t.Fatal(err)
	}

	addr, err := spout.New(contract, priv.Public())
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetQueryer(QueryHandler(querycfg))

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

	spoutcore.BeginTx(nil)
	spoutcore.Dispatcher = bolt.GetCaller("assign", token.AssignHandler(tokencfg))

	_, err = tokenSpout.Assign(addr, assignt1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.BeginTx(nil)
	spoutcore.Dispatcher = bolt.GetCaller("redeem1", RedeemHandler(cfg))

	_, err = spout.Redeem(addr, []byte(addr1), big.NewInt(0), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetQueryer(token.TokenQueryHandler(tokenQuerycfg))

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
	spoutcore.Dispatcher = bolt.GetCaller("redeem2", RedeemHandler(cfg))
	_, err = spout.Redeem(addr, []byte(addr2), bal, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetQueryer(token.TokenQueryHandler(tokenQuerycfg))
	err, data2 := tokenSpout.Account([]byte(addr2))
	if err != nil {
		t.Fatal(err)
	}

	if bal.Cmp(data2.Balance) != 0 {
		t.Fatalf("wrong redeem amount for addr2")
	}
}
