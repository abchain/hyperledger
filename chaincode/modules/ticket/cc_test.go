package ticket

import (
	"fmt"
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

var  georgeOrwell *ecdsa.PrivateKey
var  GeorgeOrwell *tx.Address

func init() {
	georgeOrwell, _ = ecdsa.NewPrivatekey(ecdsa.SECP256K1)
	GeorgeOrwell, _ = tx.NewAddress(georgeOrwell.Public())
}

func initCond() {

	cfg := NewConfig(test_tag)
	querycfg := NewConfig(test_tag)
	querycfg.SetReadOnly(true)

	cc := GeneralInvokingTemplate(test_ccname, cfg).MustMerge(GeneralQueryTemplate(test_ccname, querycfg))
	bolt = rpc.NewLocalChaincode(cc)
}

func TestAddTicket(t *testing.T) {
	initCond()
	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	ownerAddr := []byte(GeorgeOrwell.ToString())
	id := []byte("idx1")
	id2 := []byte("idx2")
	id3 := []byte("idx3")
	ticketCatalog := 10
	ticketDesc := "my ticket"

	addTicket(t, spoutcore, spout, ownerAddr, id, ticketCatalog, ticketDesc)
	addTicket(t, spoutcore, spout, ownerAddr, id2, ticketCatalog, ticketDesc)


	spoutcore.BeginTx([]byte("16666"))
	err := spout.Add(ownerAddr, id3, 0, ticketDesc)

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal("Invalid ticket")
	} else {
		t.Log(err)
	}

}


func TestAddDuplicatedTicket(t *testing.T) {
	initCond()
	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	ownerAddr := []byte(GeorgeOrwell.ToString())
	id := []byte("idx1")
	ticketCatalog := 10
	ticketDesc := "my ticket"

	addTicket(t, spoutcore, spout, ownerAddr, id, ticketCatalog, ticketDesc)

	spoutcore.BeginTx([]byte("16666"))
	err := spout.Add(ownerAddr, id, ticketCatalog, ticketDesc)

	if err != nil {
		t.Fatal(err)
	}
	_, err = spoutcore.Result().TxID()

	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func verifyTicket(spout *GeneralCall,
	owner []byte, id []byte,
	ticketCatalog int,
	ticketDesc string) bool {

	catalog, desc := spout.Query(owner, id)

	if catalog != ticketCatalog {
		return false
	}

	if desc != ticketDesc {
		return false
	}

	return true
}

func addTicket(t *testing.T, spoutcore *txgen.TxGenerator, spout *GeneralCall,
	owner []byte, id []byte, ticketCatalog int, ticketDesc string) {

	spoutcore.BeginTx([]byte("16666"))

	err := spout.Add(owner, id, ticketCatalog, ticketDesc)

	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	if !verifyTicket(spout, owner, id, ticketCatalog, ticketDesc) {
		t.Fatal(fmt.Errorf("Invalid ticket"))
	}
}


func TestRmNoneExistTicket(t *testing.T) {

	initCond()

	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	ownerAddr := []byte(GeorgeOrwell.ToString())
	id := []byte("idx1")

	spoutcore.BeginTx([]byte("16666"))

	err := spout.Apply(ownerAddr, id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}


func TestQueryNoneExistTicket(t *testing.T) {

	initCond()

	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	ownerAddr := []byte(GeorgeOrwell.ToString())
	id := []byte("idx1")

	catalog, _ := spout.Query(ownerAddr, id)

	if catalog > 0 {
		t.Fatal("Invalid ticket")
	}
	t.Log("Ticket does not exist")
}


func TestRmTicket(t *testing.T) {

	initCond()

	spoutcore := txgen.SimpleTxGen(test_ccname)
	spout := &GeneralCall{spoutcore}
	spoutcore.Dispatcher = bolt

	ownerAddr := []byte(GeorgeOrwell.ToString())
	id := []byte("idx1")
	ticketCatalog := 16
	ticketDesc := "my ticket"

	addTicket(t, spoutcore, spout, ownerAddr, id, ticketCatalog, ticketDesc)

	rmTicket(t, spoutcore, spout, ownerAddr, id)

	if verifyTicket(spout, ownerAddr, id, ticketCatalog, ticketDesc) {
		t.Fatal(fmt.Errorf("Failed to remove the ticket"))
	}
}


func rmTicket(t *testing.T, spoutcore *txgen.TxGenerator, spout *GeneralCall, owner []byte, id []byte) {

	spoutcore.BeginTx([]byte("16666"))

	err := spout.Apply(owner, id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spoutcore.Result().TxID()
	if err != nil {
		t.Fatal(err)
	}

	catalog, _ := spout.Query(owner, id)

	if catalog != -1 {
		t.Fatal(fmt.Errorf("Failed to remove the ticket"))
	}
}
