package crypto

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func GetSignatureTestData(t *testing.T) ([]byte, *PublicKey, *Signature) {

	// Generate private key

	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	// Generate message data

	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		t.Fatal("Genereate rand bytes fail: ", err)
	}

	// Signing

	sig, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
	if err != nil {
		t.Fatal("Sign raw data fail: ", err)
	}

	pub, err := priv.ChildPublic(big.NewInt(0))
	if err != nil {
		t.Fatal("Get root public key fail: ", err)
	}

	return rb, pub, sig
}

func TestSignature_Verify(t *testing.T) {

	rb, pub, sig := GetSignatureTestData(t)

	// Verify

	if !sig.Verify(rb, pub) {
		t.Fatal("Verify signature fail")
	}
}

func TestSignature_Serialize(t *testing.T) {

	rb, pub, sig := GetSignatureTestData(t)

	raw := sig.Serialize()
	t.Logf("Dump signature: %v", raw)

	sig2, err := SignatureFromBytes(raw)
	if err != nil {
		t.Fatal("Import signature fail: ", err)
	}

	if !sig.IsEqual(sig2) {
		t.Fatal("Signature not equal: ", err)
	}

	if !sig.Verify(rb, pub) {
		t.Fatal("Verify signature fail")
	}

	if !sig2.Verify(rb, pub) {
		t.Fatal("Verify signature fail")
	}
}
