package registrar

import (
	"bytes"
	"github.com/abchain/fabric/core/chaincode/shim"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/crypto"
	tx "hyperledger.abchain.org/tx"
	"math/big"
	"testing"
)

const (
	test_tag    = "test"
	test_ccname = "testCC"

	managattrN  = "role"
	regionattrN = "region"

	totalToken = "10000000000000000000000000000"
	assign1    = "500000000000000000000000000"
	assign2    = "100000000000000000000000000"
	trans1     = "400000000000000000000000000" //privk -> privk not reg
	trans2     = "200000000000000000000000000" //privk not reg -> privk
	result1    = "100000000000000000000000000"
)

var spout *GeneralCall
var bolt *rpc.DummyCallerBuilder
var stub *shim.MockStub

var cfg = &StandardRegistrarConfig{test_tag, false, managattrN, regionattrN}
var querycfg = &StandardRegistrarConfig{test_tag, true, managattrN, regionattrN}
var tokencfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}
var tokenQuerycfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}

var privkey *crypto.PrivateKey
var privkeyNotReg *crypto.PrivateKey

func assign(t *testing.T) {
	stub = shim.NewMockStub("RegTest", nil)
	spout = &GeneralCall{txgen.SimpleTxGen(test_ccname)}
	bolt = &rpc.DummyCallerBuilder{test_ccname, stub}

	tokenSpout := &token.GeneralCall{txgen.SimpleTxGen(test_ccname)}

	stub.MockTransactionStart("deployment")

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	deployargs, err := token.CCDeploy(total, nil)

	if err != nil {
		t.Fatal(err)
	}

	handlers := map[string]rpc.DeployHandler{
		token.DeployMethod: token.CCDeployHandler(test_tag)}

	err = rpc.DeployCC(stub, deployargs, handlers)

	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)

	tokenSpout.Dispatcher = bolt.GetCaller(token.AssignHandler(tokencfg))

	privkey, err = crypto.NewPrivatekey(crypto.DefaultCurveType)
	if err != nil {
		t.Fatal(err)
	}

	privkeyNotReg, err = crypto.NewPrivatekey(crypto.DefaultCurveType)
	if err != nil {
		t.Fatal(err)
	}

	addr, err := tx.NewAddressFromPrivateKey(privkey)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("assigment1")

	assignt1, ok := big.NewInt(0).SetString(assign1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	_, err = tokenSpout.Assign(addr.Hash, assignt1)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)
	addr, err = tx.NewAddressFromPrivateKey(privkeyNotReg)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("assigment2")

	assignt2, ok := big.NewInt(0).SetString(assign2, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	_, err = tokenSpout.Assign(addr.Hash, assignt2)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)
	return
}

func TestReg(t *testing.T) {
	assign(t)

	spout.Dispatcher = bolt.GetCaller(RegistrarHandler(cfg))

	stub.MockTransactionStart("registrar1")

	qkey, err := spout.Registrar(privkey.Public(), "Yosemite")
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)

	subk, err := privkey.ChildKey(big.NewInt(184467442737))
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("registrar2")

	_, err = spout.Registrar(subk.Public(), "Yosemite")
	if err == nil {
		t.Fatal("reg a childkey")
	}

	stub.MockTransactionEnd("", err)
	spout.Dispatcher = bolt.GetCaller(QueryPkHandler(querycfg))

	err, data := spout.Pubkey(qkey)
	if err != nil {
		t.Fatal(err)
	}

	if data.RegTxid != "registrar1" {
		t.Fatal("wrong txid")
	}

	qpk, err := crypto.PublicKeyFromPBMessage(data.Pk)
	if err != nil {
		t.Fatal(err)
	}

	if !qpk.IsEqual(privkey.Public()) {
		t.Fatal("unmatch pk record")
	}

	if data.Enabled {
		t.Fatal("wrong enable status")
	}

	if bytes.Compare(qkey, privkey.Public().RootFingerPrint) != 0 {
		t.Fatal("wrong pkey index")
	}

	stub.MockTransactionStart("active")

	spout.Dispatcher = bolt.GetCaller(ActivePkHandler(cfg))

	err = spout.ActivePk(qkey)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)
	spout.Dispatcher = bolt.GetCaller(QueryPkHandler(querycfg))

	err, data = spout.Pubkey(qkey)
	if err != nil {
		t.Fatal(err)
	}

	if !data.Enabled {
		t.Fatal("wrong enable status")
	}

}

func TestDirectReg(t *testing.T) {
	assign(t)

	spout.Dispatcher = bolt.GetCaller(AdminRegistrarHandler(cfg))

	stub.MockTransactionStart("adminreg1")

	err := spout.AdminRegistrar(privkey.Public())
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)
	spout.Dispatcher = bolt.GetCaller(QueryPkHandler(querycfg))

	err, data := spout.Pubkey(privkey.Public().RootFingerPrint)
	if err != nil {
		t.Fatal(err)
	}

	if data.RegTxid != "adminreg1" {
		t.Fatal("wrong txid")
	}

	if !data.Enabled {
		t.Fatal("wrong enable status")
	}

	qpk, err := crypto.PublicKeyFromPBMessage(data.Pk)
	if err != nil {
		t.Fatal(err)
	}

	if !qpk.IsEqual(privkey.Public()) {
		t.Fatal("unmatch pk record")
	}

}

func TestFund(t *testing.T) {
	TestDirectReg(t)

	tokenSpout := &token.GeneralCall{txgen.DefaultTxGen(test_ccname, privkey)}
	h := token.TransferHandler(tokencfg)
	caller := bolt.GetCaller(h)
	tokenSpout.Dispatcher = caller
	err := bolt.AppendPreHandler(caller, txhandle.AddrCredVerifier{h})
	if err != nil {
		t.Fatal(err)
	}

	err = bolt.AppendPreHandler(caller, RegistrarPreHandler(querycfg, h))
	if err != nil {
		t.Fatal(err)
	}

	transt1, ok := big.NewInt(0).SetString(trans1, 10)
	if !ok {
		t.Fatal("parse int fail")
	}

	transt2, ok := big.NewInt(0).SetString(trans1, 10)
	if !ok {
		t.Fatal("parse int fail")
	}

	addr1, err := tx.NewAddressFromPrivateKey(privkey)
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := tx.NewAddressFromPrivateKey(privkeyNotReg)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionStart("transfer1")

	_, err = tokenSpout.Transfer(addr1.Hash, addr2.Hash, transt1)
	if err != nil {
		t.Fatal(err)
	}

	stub.MockTransactionEnd("", err)
	stub.MockTransactionStart("transfer2")
	tokenSpout.Credgenerator = txgen.NewSingleKeyCred(privkeyNotReg)

	_, err = tokenSpout.Transfer(addr2.Hash, addr1.Hash, transt2)
	if err == nil {
		t.Fatal("Do transfer without reg publickey")
	}

	stub.MockTransactionEnd("", err)
	tokenSpout.Dispatcher = bolt.GetCaller(token.TokenQueryHandler(tokenQuerycfg))

	addr1bal, ok := big.NewInt(0).SetString(result1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	err, addr1data := tokenSpout.Account(addr1.Hash)
	if addr1data == nil {
		t.Fatal("Get addr data fail", err)
	}

	if addr1bal.Cmp(big.NewInt(0).SetBytes(addr1data.Balance)) != 0 {
		t.Fatal("Wrong balance for addr1")
	}

}
