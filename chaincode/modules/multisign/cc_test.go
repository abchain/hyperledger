package multisign

import (
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	tx "hyperledger.abchain.org/core/tx"
	"testing"

	"hyperledger.abchain.org/core/crypto/ecdsa"
)

const (
	test_tag    = "test"
	test_ccname = "testCC"
)

var bolt *rpc.ChaincodeAdapter

var nineteenEightyFour, georgeOrwell, animalFarm, bigBrother *ecdsa.PrivateKey
var NineteenEightyFour, GeorgeOrwell, AnimalFarm, BigBrother *tx.Address

func init() {

	nineteenEightyFour, _ = ecdsa.NewPrivatekey(ecdsa.SECP256K1)
	georgeOrwell, _ = ecdsa.NewPrivatekey(ecdsa.SECP256K1)
	animalFarm, _ = ecdsa.NewPrivatekey(ecdsa.SECP256K1)
	bigBrother, _ = ecdsa.NewPrivatekey(ecdsa.SECP256K1)

	NineteenEightyFour, _ = tx.NewAddress(nineteenEightyFour.Public())
	GeorgeOrwell, _ = tx.NewAddress(georgeOrwell.Public())
	AnimalFarm, _ = tx.NewAddress(animalFarm.Public())
	BigBrother, _ = tx.NewAddress(bigBrother.Public())

}

func initCond() {

	cfg := NewConfig(test_tag)
	querycfg := NewConfig(test_tag)
	querycfg.SetReadOnly(true)

	cc := GeneralInvokingTemplate(test_ccname, cfg).MustMerge(GeneralQueryTemplate(test_ccname, querycfg))
	bolt = rpc.NewLocalChaincode(cc)
}
func TestChaincode(t *testing.T) {

	initCond()

	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	spoutcore.BeginTx([]byte("1984"))
	addrhash, err := spout.Contract(100, map[string]int32{
		GeorgeOrwell.ToString():       70,
		NineteenEightyFour.ToString(): 60,
		AnimalFarm.ToString():         50})

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	addrs := tx.NewAddressFromHash(addrhash).ToString()
	t.Log(addrs)

	err, ret := spout.Query(addrs)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret.Addrs) != 3 {
		t.Fatal("Wrong contract")
	}

	if ind := ret.Find(GeorgeOrwell.ToString()); ind < 0 {
		t.Fatal("Not found")
	} else if ret.Addrs[ind].Weight != 70 {
		t.Fatal("Wrong record")
	}

	//check duplicated
	spoutcore.BeginTx([]byte("1984"))
	_, err = spout.Contract(100, map[string]int32{
		NineteenEightyFour.ToString(): 60,
		GeorgeOrwell.ToString():       70,
		AnimalFarm.ToString():         50})

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal("duplicated allowed")
	}

	spoutcore.BeginTx(nil)
	err = spout.Update(addrs, NineteenEightyFour.ToString(), BigBrother.ToString())

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	spoutcore.BeginTx(nil)
	err = spout.Update(addrs, AnimalFarm.ToString(), GeorgeOrwell.ToString())

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal("update to duplicated contract addr")
	}

	spoutcore.BeginTx(nil)
	err = spout.Update(addrs, AnimalFarm.ToString(), "")

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	err, ret = spout.Query(addrs)
	if err != nil {
		t.Fatal(err)
	}

	if len(ret.Addrs) != 2 {
		t.Fatal("Wrong contract")
	}

	if ind := ret.Find(BigBrother.ToString()); ind < 0 {
		t.Fatal("Not found")
	} else if ret.Addrs[ind].Weight != 60 {
		t.Fatal("Wrong record")
	}

	if ind := ret.Find(AnimalFarm.ToString()); ind >= 0 {
		t.Fatal("found removed item")
	}

}
