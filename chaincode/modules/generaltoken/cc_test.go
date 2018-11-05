package generaltoken

import (
	"bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
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
var tokencfg = &StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}
var tokenQuerycfg = &StandardTokenConfig{nonce.StandardNonceConfig{test_tag, true}}

func TestDeployCc(t *testing.T) {

	bolt.Reset()

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	deployer := &txgen.BatchTxCall{}
	deployer.TxGenerator = spoutcore
	spout := &GeneralCall{deployer}
	err := spout.Init(total)
	if err != nil {
		t.Fatal(err)
	}

	//init bolt and do deploy tx
	spoutcore.BeginDeploy(nil)
	spoutcore.Dispatcher = bolt.GetCaller("deployment",
		txhandle.BatchTxHandler(map[string]*txhandle.ChaincodeTx{
			Method_Init: &txhandle.ChaincodeTx{test_ccname, InitHandler(tokencfg), nil, nil}}))

	err = deployer.CommitBatch("init")

	if err != nil {
		t.Fatal(err)
	}

	if len(bolt.Stub().State) != 1 {
		t.Fatal("Invalid state count")
	}

	spout = &GeneralCall{spoutcore}
	//init bolt and do first query tx
	spoutcore.Dispatcher = bolt.GetQueryer(GlobalQueryHandler(tokenQuerycfg))

	err, data := spout.Global()
	if err != nil {
		t.Fatal(err)
	}

	if data.TotalTokens == nil {
		t.Fatal("Fail deploy data")
	}

	if total.Cmp(big.NewInt(0).SetBytes(data.TotalTokens)) != 0 {
		t.Fatal("Invalid amount")
	}
}

func TestAssignCc(t *testing.T) {

	TestDeployCc(t)
	spout := &GeneralCall{spoutcore}

	assignt1, ok := big.NewInt(0).SetString(assign1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	assignt2, ok := big.NewInt(0).SetString(assign2, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	assignt3, ok := big.NewInt(0).SetString(assign3, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	spoutcore.Dispatcher = bolt.GetCaller("assigment1", AssignHandler(tokencfg))
	spoutcore.BeginTx(nil)

	nc1, err := spout.Assign([]byte(addr1), assignt1)
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetCaller("assigment2", AssignHandler(tokencfg))
	spoutcore.BeginTx(nil)

	nc2, err := spout.Assign([]byte(addr2), assignt2)
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetCaller("assigment3", AssignHandler(tokencfg))
	fixednc := "fixednc"
	spoutcore.BeginTx([]byte(fixednc))

	nc3, err := spout.Assign([]byte(addr3), assignt3)
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetCaller("assigment3_fail", AssignHandler(tokencfg))
	spoutcore.BeginTx([]byte(fixednc))
	_, err = spout.Assign([]byte(addr3), assignt3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal("Execute duplicated assigment")
	}

	spoutcore.Dispatcher = bolt.GetQueryer(nonce.NonceQueryHandler(&tokenQuerycfg.StandardNonceConfig))

	err, nc1data := spout.Nonce(nc1)
	if nc1data == nil {
		t.Fatal("Get nonce data fail", err)
	}

	if assignt1.Cmp(big.NewInt(0).SetBytes(nc1data.Amount)) != 0 {
		t.Fatal("Wrong nonce data")
	}

	err, nc2data := spout.Nonce(nc2)
	if nc1data == nil {
		t.Fatal("Get nonce data fail", err)
	}

	if assignt2.Cmp(big.NewInt(0).SetBytes(nc2data.Amount)) != 0 {
		t.Fatal("Wrong nonce data")
	}

	if bytes.Compare(nonce.GeneralTokenNonceKey([]byte(fixednc), nil,
		[]byte(addr3), assignt3.Bytes()), nc3) != 0 {
		t.Fatal("Get nonce key fail")
	}

	spoutcore.Dispatcher = bolt.GetQueryer(GlobalQueryHandler(tokenQuerycfg))

	err, gdata := spout.Global()
	if err != nil {
		t.Fatal(err)
	}

	if gdata.TotalTokens == nil {
		t.Fatal("Fail deploy data")
	}

	restt, ok := big.NewInt(0).SetString(rest, 10)
	if !ok {
		t.Fatal("parse int fail")
	}

	if restt.Cmp(big.NewInt(0).SetBytes(gdata.UnassignedTokens)) != 0 {
		t.Fatal("Invalid amount")
	}
}

func TestTransferCc(t *testing.T) {

	TestAssignCc(t)
	spout := &GeneralCall{spoutcore}

	transt1, ok := big.NewInt(0).SetString(trans1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	transt2, ok := big.NewInt(0).SetString(trans2, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	transt3, ok := big.NewInt(0).SetString(trans3, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	spoutcore.Dispatcher = bolt.GetCaller("transfer1", TransferHandler(tokencfg))
	spoutcore.BeginTx(nil)

	nc1, err := spout.Transfer([]byte(addr1), []byte(addr4), transt1)
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetCaller("transfer2", TransferHandler(tokencfg))
	spoutcore.BeginTx(nil)

	nc2, err := spout.Transfer([]byte(addr2), []byte(addr4), transt2)
	if err != nil {
	}

	spoutcore.Dispatcher = bolt.GetCaller("transfer3", TransferHandler(tokencfg))
	fixednc := "fixednc"
	spoutcore.BeginTx([]byte(fixednc))

	nc3, err := spout.Transfer([]byte(addr3), []byte(addr1), transt3)
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.Dispatcher = bolt.GetCaller("transfer3_fail", TransferHandler(tokencfg))

	spoutcore.BeginTx([]byte(fixednc))
	_, err = spout.Transfer([]byte(addr3), []byte(addr1), transt3)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()

	if err == nil {
		t.Fatal("Execute duplicated transfer")
	}

	spoutcore.Dispatcher = bolt.GetCaller("transfer3_fail2", TransferHandler(tokencfg))

	transtf, ok := big.NewInt(0).SetString(transfail, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	spoutcore.BeginTx(nil)
	_, err = spout.Transfer([]byte(addr1), []byte(addr1), transtf)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()

	if err == nil {
		t.Fatal("Execute overflow transfer")
	}

	spoutcore.Dispatcher = bolt.GetQueryer(nonce.NonceQueryHandler(&tokenQuerycfg.StandardNonceConfig))

	err, nc1data := spout.Nonce(nc1)
	if nc1data == nil {
		t.Fatal("Get nonce data fail", err)
	}

	if transt1.Cmp(big.NewInt(0).SetBytes(nc1data.Amount)) != 0 {
		t.Fatal("Wrong nonce data")
	}

	err, nc2data := spout.Nonce(nc2)
	if nc1data == nil {
		t.Fatal("Get nonce data fail", err)
	}

	if transt2.Cmp(big.NewInt(0).SetBytes(nc2data.Amount)) != 0 {
		t.Fatal("Wrong nonce data")
	}

	if bytes.Compare(nonce.GeneralTokenNonceKey([]byte(fixednc), []byte(addr3),
		[]byte(addr1), transt3.Bytes()), nc3) != 0 {
		t.Fatal("Get nonce key fail")
	}

	spoutcore.Dispatcher = bolt.GetQueryer(TokenQueryHandler(tokenQuerycfg))

	addr1bal, ok := big.NewInt(0).SetString(result1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	addr4bal, ok := big.NewInt(0).SetString(result2, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	err, addr1data := spout.Account([]byte(addr1))
	if addr1data == nil {
		t.Fatal("Get addr data fail", err)
	}

	if addr1bal.Cmp(big.NewInt(0).SetBytes(addr1data.Balance)) != 0 {
		t.Fatal("Wrong balance for addr1")
	}

	err, addr4data := spout.Account([]byte(addr4))
	if addr4data == nil {
		t.Fatal("Get addr data fail", err)
	}

	if addr4bal.Cmp(big.NewInt(0).SetBytes(addr4data.Balance)) != 0 {
		t.Fatal("Wrong balance for addr4")
	}

}
