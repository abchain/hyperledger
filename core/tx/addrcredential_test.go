package abchainTx

import (
	"crypto/rand"
	"math/big"
	"testing"

	abcrypto "hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
)

var privkey1 *abcrypto.PrivateKey
var privkey2 *abcrypto.PrivateKey
var privkey3 *abcrypto.PrivateKey
var ccAddr1 *Address
var ccAddr2 *Address
var msghash []byte

func tinit(t *testing.T) {

	var err error

	privkey1, err = abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail:", err)
	}

	privkey2, err = privkey1.ChildKey(big.NewInt(184467442737))
	if err != nil {
		t.Fatal("Generate child private key fail: %v", err)
	}

	privkey3, err = abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail:", err)
	}
	// Generate message data

	msghash = make([]byte, 32)
	_, err = rand.Read(msghash)
	if err != nil {
		t.Fatal("Genereate rand bytes fail:", err)
	}

	ccAddrhash := make([]byte, 32)
	_, err = rand.Read(ccAddrhash)
	if err != nil {
		t.Fatal("Genereate rand bytes fail:", err)
	}

	ccAddr1 = NewAddressFromHash(ccAddrhash)

	ccAddrhash[0]++

	ccAddr2 = NewAddressFromHash(ccAddrhash)
}

func TestSoleCred(t *testing.T) {

	tinit(t)

	pk1 := privkey1.Public()
	pk2 := privkey2.Public()

	addr1, err := NewAddress(pk1)
	if err != nil {
		t.Fatal("Generate addr1 fail: ", err)
	}

	addr2, err := NewAddress(pk2)
	if err != nil {
		t.Fatal("Generate addr2 fail: ", err)
	}

	sig1, err := privkey1.SignwithThis(rand.Reader, msghash)
	if err != nil {
		t.Fatal("Generate sig1 fail: ", err)
	}

	builder := NewAddrCredentialBuilder()

	builder.AddSignature(pk1, &abcrypto.ECSignature{*sig1})

	msg := &pb.TxCredential{}

	err = builder.Update(msg)

	if err != nil {
		t.Fatal("Create builder fail: ", err)
	}

	cred, err := NewAddrCredential(msghash, msg.Addrc)

	if err != nil {
		t.Fatal("Create verifier fail: ", err)
	}

	if cred.GetCredPubkey(*addr2) != nil {
		t.Fatal("Get cred pk from wrong addr ")
	}

	pkget := cred.GetCredPubkey(*addr1)

	if pkget == nil {
		t.Fatal("Get cred pk fail ")
	}

	if !pkget.IsEqual(pk1) {
		t.Fatal("Obtain pk is not identical ")
	}

	err = cred.Verify(*addr1)
	if err != nil {
		t.Fatal("verify addr1 fail: ", err)
	}

	err = cred.Verify(*addr2)
	if err == nil {
		t.Fatal("error verify for addr2")
	}

}

func TestMutipleCred(t *testing.T) {

	tinit(t)

	pk1 := privkey1.Public()
	pk2 := privkey2.Public()
	pk3 := privkey3.Public()

	addr1, err := NewAddress(pk1)
	if err != nil {
		t.Fatal("Generate addr1 fail: ", err)
	}

	addr2, err := NewAddress(pk2)
	if err != nil {
		t.Fatal("Generate addr2 fail: ", err)
	}

	addr3, err := NewAddress(pk3)
	if err != nil {
		t.Fatal("Generate addr3 fail: ", err)
	}

	sig1, err := privkey1.SignwithThis(rand.Reader, msghash)
	if err != nil {
		t.Fatal("Generate sig1 fail: ", err)
	}

	sig3, err := privkey3.SignwithThis(rand.Reader, msghash)
	if err != nil {
		t.Fatal("Generate sig3 fail: ", err)
	}

	builder := NewAddrCredentialBuilder()

	builder.AddSignature(pk1, &abcrypto.ECSignature{*sig1})
	builder.AddSignature(pk3, &abcrypto.ECSignature{*sig3})

	msg := &pb.TxCredential{}

	err = builder.Update(msg)

	if err != nil {
		t.Fatal("Create builder fail: ", err)
	}

	cred, err := NewAddrCredential(msghash, msg.Addrc)

	if err != nil {
		t.Fatal("Create verifier fail: ", err)
	}

	err = cred.Verify(*addr1)
	if err != nil {
		t.Fatal("verify addr1 fail: ", err)
	}

	err = cred.Verify(*addr2)
	if err == nil {
		t.Fatal("error verify for addr2")
	}

	err = cred.Verify(*addr3)
	if err != nil {
		t.Fatal("verify addr3 fail: ", err)
	}

	if cred.GetCredPubkey(*addr2) != nil {
		t.Fatal("Get cred pk from wrong addr ")
	}

	pkget := cred.GetCredPubkey(*addr1)

	if pkget == nil {
		t.Fatal("Get cred pk fail ")
	}

	if !pkget.IsEqual(pk1) {
		t.Fatal("Obtain pk is not identical ")
	}

}

func TestCcCred(t *testing.T) {
	tinit(t)

	pk1 := privkey1.Public()

	var err error
	//	addr1, err := NewAddress(pk1)
	// if err != nil {
	// 	t.Fatal("Generate addr1 fail: ", err)
	// }

	builder := NewAddrCredentialBuilder()

	builder.AddCc("test1", *ccAddr1, pk1)

	msg := &pb.TxCredential{}

	err = builder.Update(msg)

	if err == nil {
		t.Fatal("Expected cc credential fail but passed")
	}

	// cred, err := NewAddrCredential(msghash, msg.Addrc)

	// if err != nil {
	// 	t.Fatal("Create verifier fail: ", err)
	// }

	// err = cred.Verify(*ccAddr1)
	// if err != nil {
	// 	//Notice: verify should be wrong
	// 	t.Log("verify addr1 fail: ", err)
	// }

	// err = cred.Verify(*addr1)
	// if err == nil {
	// 	t.Fatal("error verify for pk addr1")
	// }

	// pkget := cred.GetCredPubkey(*ccAddr1)

	// if pkget == nil {
	// 	t.Fatal("Get cred pk fail ")
	// }

	// if !pkget.IsEqual(pk1) {
	// 	t.Fatal("Obtain pk is not identical ")
	// }

	// if cred.GetCredPubkey(*addr1) != nil {
	// 	t.Fatal("Get cred pk from wrong addr ")
	// }
}

func TestMixedCred(t *testing.T) {

}
