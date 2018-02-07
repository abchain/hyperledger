package generaltoken

import (
	"bytes"
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
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
var bolt *rpc.DummyCallerBuilder
var stub *shim.MockStub
var tokencfg = &StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}
var tokenQuerycfg = &StandardTokenConfig{nonce.StandardNonceConfig{test_tag, true}}

func TestDeployCc(t *testing.T) {

	//we only use the db in stub and never do mocking from chaincode interface
	stub = shim.NewMockStub("TokenTest", nil)

	spout = &GeneralCall{txgen.SimpleTxGen(test_ccname)}
	bolt = &rpc.DummyCallerBuilder{test_ccname, stub}

	stub.MockTransactionStart("deployment")

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	deployargs, err := CCDeploy(total, nil)

	if err != nil {
		t.Fatal(err)
	}

	var h CCDeployHandler = test_tag
	handlers := make(map[string]rpc.DeployHandler)
	handlers[DeployMethod] = h

	err = rpc.DeployCC(stub, deployargs, handlers)

	if err != nil {
		t.Fatal(err)
	}

	if len(stub.State) != 1 {
		t.Fatal("Invalid state count")
	}

	//init bolt and do first query tx
	spout.Dispatcher = bolt.GetCaller(GlobalQueryHandler(tokenQuerycfg))

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

	spout.Dispatcher = bolt.GetCaller(AssignHandler(tokencfg))

	stub.MockTransactionStart("assigment1")

	nc1, err := spout.Assign([]byte(addr1), assignt1)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("assigment2")

	nc2, err := spout.Assign([]byte(addr2), assignt2)
	if err != nil {
	}

	stub.MockTransactionStart("assigment3")
	fixednc := "fixednc"
	spout.BeginTx([]byte(fixednc))

	nc3, err := spout.Assign([]byte(addr3), assignt3)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("assigment3_fail")

	spout.BeginTx([]byte(fixednc))
	_, err = spout.Assign([]byte(addr3), assignt3)
	if err == nil {
		t.Fatal("Execute duplicated assigment")
	}

	spout.Dispatcher = bolt.GetCaller(nonce.NonceQueryHandler(&tokenQuerycfg.StandardNonceConfig))

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

	spout.Dispatcher = bolt.GetCaller(GlobalQueryHandler(tokenQuerycfg))

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

	spout.Dispatcher = bolt.GetCaller(TransferHandler(tokencfg))

	stub.MockTransactionStart("transfer1")

	nc1, err := spout.Transfer([]byte(addr1), []byte(addr4), transt1)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("transfer2")

	nc2, err := spout.Transfer([]byte(addr2), []byte(addr4), transt2)
	if err != nil {
	}

	stub.MockTransactionStart("transfer3")
	fixednc := "fixednc"
	spout.BeginTx([]byte(fixednc))

	nc3, err := spout.Transfer([]byte(addr3), []byte(addr1), transt3)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("transfer3_fail")

	spout.BeginTx([]byte(fixednc))
	_, err = spout.Transfer([]byte(addr3), []byte(addr1), transt3)
	if err == nil {
		t.Fatal("Execute duplicated transfer")
	}

	stub.MockTransactionStart("transfer3_fail")

	transtf, ok := big.NewInt(0).SetString(transfail, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	spout.BeginTx(nil)
	_, err = spout.Transfer([]byte(addr1), []byte(addr1), transtf)
	if err == nil {
		t.Fatal("Execute overflow transfer")
	}

	spout.Dispatcher = bolt.GetCaller(nonce.NonceQueryHandler(&tokenQuerycfg.StandardNonceConfig))

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

	spout.Dispatcher = bolt.GetCaller(TokenQueryHandler(tokenQuerycfg))

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
