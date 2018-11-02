package registrar

import (
	"bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	"hyperledger.abchain.org/core/crypto"
	tx "hyperledger.abchain.org/core/tx"
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
var bolt = &rpc.DummyCallerBuilder{CCName: test_ccname}

var cfg = &StandardRegistrarConfig{test_tag, false, managattrN, regionattrN}
var querycfg = &StandardRegistrarConfig{test_tag, true, managattrN, regionattrN}
var tokencfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}
var tokenQuerycfg = &token.StandardTokenConfig{nonce.StandardNonceConfig{test_tag, false}}

var privkey *crypto.PrivateKey
var privkeyNotReg *crypto.PrivateKey

func assign(t *testing.T) {
	bolt.Reset()

	spout = &GeneralCall{txgen.SimpleTxGen(test_ccname)}
	tokenSpout := &token.GeneralCall{txgen.SimpleTxGen(test_ccname)}

	total, ok := big.NewInt(0).SetString(totalToken, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	deployTx := txgen.NewDeployTx()
	regD := &DeployCall{deployTx}
	regD.InitDebugMode()
	tokenD := token.DeployCall{deployTx}
	tokenD.Init(total)

	spout.Dispatcher = bolt.GetCaller("deployment",
		txhandle.DeployTxHandler(map[string]txhandle.TxHandler{
			DeployMethod:       CCDeployHandler(cfg),
			token.DeployMethod: token.CCDeployHandler(tokencfg),
		}))

	deployTx.TxGenerator = spout.TxGenerator

	err := deployTx.Deploy("init")

	if err != nil {
		t.Fatal(err)
	}

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

	assignt1, ok := big.NewInt(0).SetString(assign1, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	tokenSpout.Dispatcher = bolt.GetCaller("assigment1", token.AssignHandler(tokencfg))
	_, err = tokenSpout.Assign(addr.Hash, assignt1)
	if err != nil {
		t.Fatal(err)
	}

	addr, err = tx.NewAddressFromPrivateKey(privkeyNotReg)
	if err != nil {
		t.Fatal(err)
	}

	assignt2, ok := big.NewInt(0).SetString(assign2, 10)

	if !ok {
		t.Fatal("parse int fail")
	}

	tokenSpout.Dispatcher = bolt.GetCaller("assigment2", token.AssignHandler(tokencfg))
	_, err = tokenSpout.Assign(addr.Hash, assignt2)
	if err != nil {
		t.Fatal(err)
	}

	return
}

func TestReg(t *testing.T) {
	assign(t)

	spout.Dispatcher = bolt.GetCaller("registrar1", RegistrarHandler(cfg))
	qkey, err := spout.Registrar(privkey.Public(), "Yosemite")
	if err != nil {
		t.Fatal(err)
	}

	subk, err := privkey.ChildKey(big.NewInt(184467442737))
	if err != nil {
		t.Fatal(err)
	}

	spout.Dispatcher = bolt.GetCaller("registrar2", RegistrarHandler(cfg))
	_, err = spout.Registrar(subk.Public(), "Yosemite")
	if err == nil {
		t.Fatal("reg a childkey")
	}

	spout.Dispatcher = bolt.GetQueryer(QueryPkHandler(querycfg))

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

	spout.Dispatcher = bolt.GetCaller("active", ActivePkHandler(cfg))

	err = spout.ActivePk(qkey)
	if err != nil {
		t.Fatal(err)
	}

	spout.Dispatcher = bolt.GetQueryer(QueryPkHandler(querycfg))

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

	spout.Dispatcher = bolt.GetCaller("adminreg1", AdminRegistrarHandler(cfg))
	err := spout.AdminRegistrar(privkey.Public())
	if err != nil {
		t.Fatal(err)
	}

	spout.Dispatcher = bolt.GetQueryer(QueryPkHandler(querycfg))

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
	caller := bolt.GetCaller("transfer1", h)
	tokenSpout.Dispatcher = caller
	err := bolt.AppendPreHandler(txhandle.AddrCredVerifier{h, nil})
	if err != nil {
		t.Fatal(err)
	}

	err = bolt.AppendPreHandler(RegistrarPreHandler(querycfg, h))
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

	_, err = tokenSpout.Transfer(addr1.Hash, addr2.Hash, transt1)
	if err != nil {
		t.Fatal(err)
	}

	bolt.NewTxID("transfer2")
	tokenSpout.Credgenerator = txgen.NewSingleKeyCred(privkeyNotReg)

	_, err = tokenSpout.Transfer(addr2.Hash, addr1.Hash, transt2)
	if err == nil {
		t.Fatal("Do transfer without reg publickey")
	}

	tokenSpout.Dispatcher = bolt.GetQueryer(token.TokenQueryHandler(tokenQuerycfg))

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
