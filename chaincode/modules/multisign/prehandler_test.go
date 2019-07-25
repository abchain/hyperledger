package multisign

import (
	_ "bytes"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"

	tx "hyperledger.abchain.org/core/tx"
	"math/big"
	"testing"
)

const (
	totalToken = "10000000000000000000000000000"
)

var testCC txhandle.CollectiveTxs
var tokenbolt *rpc.ChaincodeAdapter
var spoutcore *txgen.TxGenerator

func initFullCond(multicc bool) {

	cfg := token.NewConfig(test_tag)
	querycfg := token.NewConfig(test_tag)
	querycfg.SetReadOnly(true)

	cc := token.GeneralAdminTemplate(test_ccname, cfg)
	cc = cc.MustMerge(token.GeneralInvokingTemplate(test_ccname, cfg))
	cc = cc.MustMerge(token.GeneralQueryTemplate(test_ccname, querycfg))

	func() {
		cfg := NewConfig(test_tag)
		querycfg := NewConfig(test_tag)
		querycfg.SetReadOnly(true)

		if testCC == nil {
			testCC = GeneralInvokingTemplate(test_ccname, cfg).MustMerge(GeneralQueryTemplate(test_ccname, querycfg))
		}

		if multicc {

			bolt = rpc.NewLocalChaincode(txhandle.CollectiveTxs_InnerSupport(testCC))
			bolt.Name = "mainCC"

			token.ExtendInvokingTemplate(cc,
				MultiSignAddrPreHandler(InnerInvokeConfig{txgen.InnerChaincode(bolt.Name)}))
			tokenbolt = rpc.NewLocalChaincode(cc)
			tokenbolt.Invokables[bolt.Name] = bolt.MockStub

		} else {
			token.ExtendInvokingTemplate(cc, MultiSignAddrPreHandler(cfg))
			cc = cc.MustMerge(testCC)
			tokenbolt = rpc.NewLocalChaincode(cc)
			bolt = tokenbolt
		}

		testCC = nil
	}()

}

func initTest(t *testing.T) {

	spoutcore = txgen.SimpleTxGen(test_ccname)
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

func afterInvoking(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}
}

func afterFailInvoking(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal("Unexpected passed tx")
	}
}

func doTestVerifier(t *testing.T) {

	//first build a contract
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	//a 1-of-2 multisign
	spoutcore.BeginTx(nil)
	addrhash, err := spout.Contract(50, map[string]int32{
		GeorgeOrwell.ToString():       50,
		NineteenEightyFour.ToString(): 50})

	afterInvoking(err, t)

	ctaddr1 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr1.ToString())

	tokenspout := &token.GeneralCall{spoutcore}
	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginTx(nil)

	_, err = tokenspout.Assign(AnimalFarm.Hash, big.NewInt(1000))

	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(animalFarm)
	_, err = tokenspout.Transfer(AnimalFarm.Hash, ctaddr1.Hash, big.NewInt(900))
	afterInvoking(err, t)

	if err, ctacc := tokenspout.Account(ctaddr1.Hash); err != nil {
		t.Fatal(err)
	} else if ctacc.Balance.Cmp(big.NewInt(900)) != 0 {
		t.Fatalf("Wrong balance: %s", ctacc.Balance)
	}

	//verify it is impossible to do transfer with wrong credential (animalFarm's)
	spoutcore.BeginTx(nil)
	_, err = tokenspout.Transfer(ctaddr1.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr1.Hash, AnimalFarm.Hash, big.NewInt(500))
	afterInvoking(err, t)

	if err, ctacc := tokenspout.Account(AnimalFarm.Hash); err != nil {
		t.Fatal(err)
	} else if ctacc.Balance.Cmp(big.NewInt(600)) != 0 {
		t.Fatalf("Wrong balance: %s", ctacc.Balance)
	}

	//a 2-of-2 multisign
	spoutcore.Dispatcher = bolt
	spoutcore.BeginTx(nil)
	addrhash, err = spout.Contract(70, map[string]int32{
		GeorgeOrwell.ToString():       50,
		NineteenEightyFour.ToString(): 50})

	afterInvoking(err, t)

	ctaddr2 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr2.ToString())

	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(animalFarm)
	_, err = tokenspout.Transfer(AnimalFarm.Hash, ctaddr2.Hash, big.NewInt(500))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewMultiKeyCred(georgeOrwell, nineteenEightyFour)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterInvoking(err, t)

	if err, ctacc := tokenspout.Account(AnimalFarm.Hash); err != nil {
		t.Fatal(err)
	} else if ctacc.Balance.Cmp(big.NewInt(500)) != 0 {
		t.Fatalf("Wrong balance: %s", ctacc.Balance)
	}
}

func TestVerifier(t *testing.T) {

	initFullCond(false)
	initTest(t)
	doTestVerifier(t)
}

func TestVerifier_Multicc(t *testing.T) {

	initFullCond(true)
	initTest(t)
	doTestVerifier(t)
}

func doRecursiveTestVerifier(t *testing.T) {

	//first build a contract
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	//a 1-of-2 multisign
	spoutcore.BeginTx(nil)
	addrhash, err := spout.Contract(50, map[string]int32{
		GeorgeOrwell.ToString():       50,
		NineteenEightyFour.ToString(): 50})

	afterInvoking(err, t)

	ctaddr1 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr1.ToString())

	//a 2-of-2 multisign, recursive the previous address
	spoutcore.Dispatcher = bolt
	spoutcore.BeginTx(nil)
	addrhash, err = spout.Contract(70, map[string]int32{
		ctaddr1.ToString():    50,
		BigBrother.ToString(): 50})

	afterInvoking(err, t)

	ctaddr2 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr2.ToString())

	tokenspout := &token.GeneralCall{spoutcore}
	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginTx(nil)
	_, err = tokenspout.Assign(AnimalFarm.Hash, big.NewInt(1000))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(animalFarm)
	_, err = tokenspout.Transfer(AnimalFarm.Hash, ctaddr2.Hash, big.NewInt(900))
	afterInvoking(err, t)

	//only george's credential is not enough
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

	//also bigbrother's
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(bigBrother)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

	//but ok for both
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewMultiKeyCred(bigBrother, georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterInvoking(err, t)

	if err, ctacc := tokenspout.Account(AnimalFarm.Hash); err != nil {
		t.Fatal(err)
	} else if ctacc.Balance.Cmp(big.NewInt(500)) != 0 {
		t.Fatalf("Wrong balance: %s", ctacc.Balance)
	}
}

func doRecursiveLimitTest(t *testing.T) {

	//first build a contract
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	//a 1-of-2 multisign
	spoutcore.BeginTx(nil)
	addrhash, err := spout.Contract(50, map[string]int32{
		GeorgeOrwell.ToString():       50,
		NineteenEightyFour.ToString(): 50})

	afterInvoking(err, t)

	ctaddr1 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr1.ToString())

	//a 2-of-2 multisign, recursive the previous address
	spoutcore.Dispatcher = bolt
	spoutcore.BeginTx(nil)
	addrhash, err = spout.Contract(70, map[string]int32{
		ctaddr1.ToString():    50,
		BigBrother.ToString(): 50})

	afterInvoking(err, t)

	ctaddr2 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr2.ToString())

	tokenspout := &token.GeneralCall{spoutcore}
	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginTx(nil)
	_, err = tokenspout.Assign(AnimalFarm.Hash, big.NewInt(1000))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(animalFarm)
	_, err = tokenspout.Transfer(AnimalFarm.Hash, ctaddr2.Hash, big.NewInt(900))
	afterInvoking(err, t)

	//can not pass for recursive limit
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewMultiKeyCred(bigBrother, georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

}

func TestRecursiveVerifier(t *testing.T) {

	defer func(d int) {
		defaultRecursiveDepth = d
	}(defaultRecursiveDepth)

	initFullCond(false)
	initTest(t)
	doRecursiveTestVerifier(t)

	defaultRecursiveDepth = 0
	initFullCond(false)
	initTest(t)
	doRecursiveLimitTest(t)
}

func TestRecursiveVerifier_Multicc(t *testing.T) {

	initFullCond(true)
	initTest(t)
	doTestVerifier(t)
}

func doRecursiveTestSpecial(t *testing.T) {

	//first build a contract
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	//a 1-of-2 multisign
	spoutcore.BeginTx(nil)
	addrhash, err := spout.Contract(50, map[string]int32{
		NineteenEightyFour.ToString(): 50,
		BigBrother.ToString():         50})

	afterInvoking(err, t)

	ctaddr1 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr1.ToString())

	//a 1 or 2-of-3 multisign, recursive the previous address
	spoutcore.Dispatcher = bolt
	spoutcore.BeginTx(nil)
	addrhash, err = spout.Contract(70, map[string]int32{
		ctaddr1.ToString():      50,
		BigBrother.ToString():   70,
		GeorgeOrwell.ToString(): 50})

	afterInvoking(err, t)

	ctaddr2 := tx.NewAddressFromHash(addrhash)
	t.Log(ctaddr2.ToString())

	tokenspout := &token.GeneralCall{spoutcore}
	spoutcore.Dispatcher = tokenbolt
	spoutcore.BeginTx(nil)
	_, err = tokenspout.Assign(AnimalFarm.Hash, big.NewInt(1000))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(animalFarm)
	_, err = tokenspout.Transfer(AnimalFarm.Hash, ctaddr2.Hash, big.NewInt(900))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(400))
	afterFailInvoking(err, t)

	//you need both george and 1984, or just bigbrother
	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewMultiKeyCred(nineteenEightyFour, georgeOrwell)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(200))
	afterInvoking(err, t)

	spoutcore.BeginTx(nil)
	spoutcore.Credgenerator = txgen.NewSingleKeyCred(bigBrother)
	_, err = tokenspout.Transfer(ctaddr2.Hash, AnimalFarm.Hash, big.NewInt(200))
	afterInvoking(err, t)

	if err, ctacc := tokenspout.Account(AnimalFarm.Hash); err != nil {
		t.Fatal(err)
	} else if ctacc.Balance.Cmp(big.NewInt(500)) != 0 {
		t.Fatalf("Wrong balance: %s", ctacc.Balance)
	}
}

func TestRecursiveSpecial(t *testing.T) {

	initFullCond(false)
	initTest(t)
	doRecursiveTestSpecial(t)

}

func TestRecursiveSpecial_Multicc(t *testing.T) {

	initFullCond(true)
	initTest(t)
	doRecursiveTestSpecial(t)

}
